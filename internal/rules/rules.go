// Package rules is used to manipulate traffic control rules.
package rules

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"time"

	"github.com/kakeetopius/qosm/internal/core/nft"
	"github.com/kakeetopius/qosm/internal/core/tc"
	"github.com/kakeetopius/qosm/internal/db"
	"github.com/kakeetopius/qosm/internal/util"
)

type Rule struct {
	ID        int
	Target    string
	Type      string
	Priority  string
	CreatedAt time.Time
}

func AddDomainRule(dbCon *sql.DB, htbCtx *tc.HTBCtx, domain string, priority string, logger *slog.Logger) (Rule, error) {
	exists, err := db.CheckDomainRuleExists(dbCon, domain)
	rule := Rule{}
	if err != nil {
		return rule, err
	}
	if exists {
		return rule, fmt.Errorf("rule for %v already exists", domain)
	}

	var prio tc.Priority
	switch priority {
	case "high":
		prio = tc.PRIORITYHIGH
	case "low":
		prio = tc.PRIORITYLOW
	default:
		return rule, fmt.Errorf("unknown priority: %s", priority)
	}

	_, err = netip.ParseAddr(domain)
	if err == nil {
		return rule, fmt.Errorf("%v seems to be an IP address not a domain", domain)
	}

	util.Debug(logger, "resolving_domain", "domain", domain)
	ips, err := net.LookupIP(domain)
	if err != nil {
		util.Error(logger, "resolve_error", "domain", domain, "error", err.Error())
		return rule, err
	}
	addrs := util.NetIPtoNetIPPRefix(ips)

	util.Debug(logger, "add_rule", "target", domain, "priority", priority)

	err = htbCtx.AddRule(addrs, prio)
	if err != nil {
		util.Error(logger, "tc_error", "error", err.Error())
		return rule, err
	}
	err = db.AddDomainToPriority(dbCon, domain, priority, addrs)
	if err != nil {
		return rule, err
	}

	domainRule, err := db.GetDomainRuleNameByWithoutIPs(dbCon, domain)
	if err != nil {
		return rule, err
	}

	return Rule{
		Type:      "domain",
		Priority:  domainRule.Priority,
		Target:    domainRule.DomainName,
		ID:        domainRule.ID,
		CreatedAt: domainRule.CreatedAt,
	}, nil
}

func AddIPRule(dbCon *sql.DB, htbCtx *tc.HTBCtx, ip string, priority string, logger *slog.Logger) (Rule, error) {
	exists, err := db.CheckIPRuleExists(dbCon, ip)
	rule := Rule{}
	if err != nil {
		return rule, err
	}
	if exists {
		return rule, fmt.Errorf("rule for %v already exists", ip)
	}

	var prio tc.Priority
	switch priority {
	case "high":
		prio = tc.PRIORITYHIGH
	case "low":
		prio = tc.PRIORITYLOW
	default:
		return rule, fmt.Errorf("unknown priority: %s", priority)
	}

	addrs, err := util.TargetsFromString(ip)
	if err != nil {
		return rule, fmt.Errorf("invalid IP address: %v", ip)
	}

	util.Debug(logger, "add_rule", "target", ip, "priority", priority)

	err = htbCtx.AddRule(addrs, prio)
	if err != nil {
		util.Error(logger, "tc_error", "error", err.Error())
		return rule, err
	}

	ipString := addrs[0].String()
	err = db.AddIPToPriority(dbCon, ipString, priority)
	if err != nil {
		return rule, err
	}

	ipRule, err := db.GetIPRuleByName(dbCon, ipString)
	if err != nil {
		return rule, err
	}

	return Rule{
		Type:      "ip",
		Priority:  ipRule.Priority,
		Target:    ipRule.IP,
		ID:        ipRule.ID,
		CreatedAt: ipRule.CreatedAt,
	}, nil
}

func DeleteDomainRuleByID(dbConn *sql.DB, htbCtx *tc.HTBCtx, domainRuleID int) error {
	domainRule, err := db.GetDomainRuleByID(dbConn, domainRuleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rules to delete for domain with ID %v", domainRuleID)
		}
		return err
	}

	err = db.DeleteDomainRuleByID(dbConn, domainRuleID, domainRule.Priority)
	if err != nil {
		return err
	}

	return deleteDomainRule(domainRule, htbCtx)
}

func DeleteDomainRuleByName(dbConn *sql.DB, htbCtx *tc.HTBCtx, name string) error {
	domainRule, err := db.GetDomainRuleByName(dbConn, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rules to delete for domain %v", name)
		}
		return err
	}

	err = db.DeleteDomainRuleByName(dbConn, name, domainRule.Priority)
	if err != nil {
		return err
	}

	return deleteDomainRule(domainRule, htbCtx)
}

func DeleteIPRuleByID(dbConn *sql.DB, htbCtx *tc.HTBCtx, ipRuleID int) error {
	ipRule, err := db.GetIPRuleByID(dbConn, ipRuleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rules to delete for IP rule with ID %v", ipRuleID)
		}
		return err
	}

	err = db.DeleteIPRuleByID(dbConn, ipRuleID, ipRule.Priority)
	if err != nil {
		return err
	}

	return deleteIPRule(htbCtx, ipRule)
}

func DeleteIPRuleByName(dbConn *sql.DB, htbCtx *tc.HTBCtx, ipRuleName string) error {
	ipRule, err := db.GetIPRuleByName(dbConn, ipRuleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rules to delete for ip %v", ipRuleName)
		}
		return err
	}

	err = db.DeleteIPRuleByName(dbConn, ipRuleName, ipRule.Priority)
	if err != nil {
		return err
	}

	return deleteIPRule(htbCtx, ipRule)
}

func DeleteAllRules(dbConn *sql.DB) error {
	err := nft.DeleteTable()
	if err != nil {
		if !errors.Is(err, nft.ErrTableNotFound) {
			return err
		}
	}

	err = db.FlushDomainRules(dbConn)
	if err != nil {
		return err
	}

	err = db.FlushIPRules(dbConn)
	if err != nil {
		return err
	}

	return nil
}

func GetAll(dbCon *sql.DB) ([]Rule, error) {
	ipRules, err := db.GetAllIPRules(dbCon)
	if err != nil {
		return nil, err
	}
	domainRules, err := db.GetAllDomainRulesWithoutIPs(dbCon)
	if err != nil {
		return nil, err
	}

	return joinIPAndDomainRules(ipRules, domainRules), nil
}

func GetHighPriority(dbCon *sql.DB) ([]Rule, error) {
	highPrioIPRules, err := db.GetHighPrioIPs(dbCon)
	if err != nil {
		return nil, err
	}
	highPrioDomainRules, err := db.GetHighPrioDomains(dbCon)
	if err != nil {
		return nil, err
	}

	return joinIPAndDomainRules(highPrioIPRules, highPrioDomainRules), nil
}

func GetLowPriority(dbCon *sql.DB) ([]Rule, error) {
	lowPrioIPRules, err := db.GetLowPrioIPs(dbCon)
	if err != nil {
		return nil, err
	}
	lowPrioDomainRules, err := db.GetLowPrioDomains(dbCon)
	if err != nil {
		return nil, err
	}

	return joinIPAndDomainRules(lowPrioIPRules, lowPrioDomainRules), nil
}

func joinIPAndDomainRules(ipRules []db.IPRule, domainRules []db.DomainRule) []Rule {
	allRules := make([]Rule, 0, len(ipRules)+len(domainRules))
	for _, rule := range ipRules {
		allRules = append(allRules, Rule{
			ID:        rule.ID,
			Priority:  rule.Priority,
			Target:    rule.IP,
			Type:      "ip",
			CreatedAt: rule.CreatedAt,
		})
	}

	for _, rule := range domainRules {
		allRules = append(allRules, Rule{
			ID:        rule.ID,
			Priority:  rule.Priority,
			Target:    rule.DomainName,
			Type:      "domain",
			CreatedAt: rule.CreatedAt,
		})
	}

	return allRules
}

func deleteDomainRule(domainRule db.DomainRule, htbCtx *tc.HTBCtx) error {
	addrs := make([]netip.Prefix, 0, len(domainRule.IPs))
	for _, addr := range domainRule.IPs {
		ip, iperr := netip.ParsePrefix(addr.IP)
		if iperr != nil {
			return iperr
		}
		addrs = append(addrs, ip)
	}

	switch domainRule.Priority {
	case "high":
		return htbCtx.NFTFilter.DeleteTargetFromHighPriority(addrs)
	case "low":
		return htbCtx.NFTFilter.DeleteTargetFromLowPriority(addrs)
	default:
		return fmt.Errorf("unknown priority: %v", domainRule.Priority)
	}
}

func deleteIPRule(htbCtx *tc.HTBCtx, ipRule db.IPRule) error {
	addr, err := netip.ParsePrefix(ipRule.IP)
	if err != nil {
		return err
	}

	switch ipRule.Priority {
	case "high":
		return htbCtx.NFTFilter.DeleteTargetFromHighPriority([]netip.Prefix{addr})
	case "low":
		return htbCtx.NFTFilter.DeleteTargetFromLowPriority([]netip.Prefix{addr})
	default:
		return fmt.Errorf("unknown priority: %v", ipRule.Priority)
	}
}
