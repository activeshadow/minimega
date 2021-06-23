package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"minicli"
	log "minilog"
)

func init() {
	minicli.Register(&minicli.Handler{
		Patterns: []string{
			"fw <allow,> <in,out> <index> <dst> <proto>",
			"fw <allow,> <in,out> <index> <src> <dst> <proto>",
			"fw <flush,>",
		},
		Call: handleFW,
	})
}

func handleFW(c *minicli.Command, r chan<- minicli.Responses) {
	defer func() {
		r <- nil
	}()

	if c.BoolArgs["flush"] {
		log.Debug("flushing firwall")

		if err := flushFW(); err != nil {
			log.Errorln(err)
		}

		return
	}

	if c.BoolArgs["allow"] {
		var (
			idx   = -1
			iface = "lo"
			err   error
		)

		if c.StringArgs["index"] != "lo" {
			idx, err = strconv.Atoi(c.StringArgs["index"])
			if err != nil {
				log.Errorln("converting interface index: %v", err)
				return
			}
		}

		if idx != -1 {
			// get interface name using the index
			if iface, err = findEth(idx); err != nil {
				log.Errorln("getting interface name for index: %v", err)
				return
			}
		}

		var rule []string

		if c.BoolArgs["in"] {
			rule = []string{"FORWARD", "-i", iface}
		} else {
			rule = []string{"FORWARD", "-o", iface}
		}

		if proto := c.StringArgs["proto"]; proto != "" {
			rule = append(rule, "-p", proto)
		}

		if src := c.StringArgs["src"]; src != "" {
			fields := strings.Split(src, ":")

			switch len(fields) {
			case 1:
				rule = append(rule, "-s", src)
			case 2:
				if _, err = strconv.Atoi(fields[1]); err != nil {
					log.Errorln("converting source port: %v", err)
					return
				}

				if proto := c.StringArgs["proto"]; proto == "" {
					log.Errorln("must specify proto when specifying ports")
					return
				}

				rule = append(rule, "-s", fields[0], "--sport", fields[1])
			default:
				log.Errorln("malformed source")
				return
			}
		}

		if dst := c.StringArgs["dst"]; dst != "" {
			fields := strings.Split(dst, ":")

			switch len(fields) {
			case 1:
				rule = append(rule, "-d", dst)
			case 2:
				if _, err = strconv.Atoi(fields[1]); err != nil {
					log.Errorln("converting destination port: %v", err)
					return
				}

				if proto := c.StringArgs["proto"]; proto == "" {
					log.Errorln("must specify proto when specifying ports")
					return
				}

				rule = append(rule, "-d", fields[0], "--dport", fields[1])
			default:
				log.Errorln("malformed destination")
				return
			}
		}

		rule = append(rule, "-j", "ACCEPT")

		if err := addRule(rule); err != nil {
			log.Errorln(err)
			return
		}

		if err := defaultForwardDrop(); err != nil {
			log.Errorln(err)
		}
	}
}

func addRule(rule []string) error {
	args := append([]string{"-A"}, rule...)

	out, err := exec.Command("iptables", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("adding iptables rule (%s) %v: %v", strings.Join(rule, " "), err, string(out))
	}

	return nil
}

func flushFW() error {
	if out, err := exec.Command("iptables", "-F").CombinedOutput(); err != nil {
		return fmt.Errorf("flushing rules %v: %v", err, string(out))
	}

	if out, err := exec.Command("iptables", "-X").CombinedOutput(); err != nil {
		return fmt.Errorf("deleting custom chains %v: %v", err, string(out))
	}

	if err := defaultAllAccept(); err != nil {
		return err
	}

	if err := createEstablishedChain(); err != nil {
		return err
	}

	return nil
}

func defaultAllAccept() error {
	if out, err := exec.Command("iptables", "-P", "INPUT", "ACCEPT").CombinedOutput(); err != nil {
		return fmt.Errorf("defaulting INPUT to ACCEPT %v: %v", err, string(out))
	}

	if out, err := exec.Command("iptables", "-P", "FORWARD", "ACCEPT").CombinedOutput(); err != nil {
		return fmt.Errorf("defaulting FORWARD to ACCEPT %v: %v", err, string(out))
	}

	if out, err := exec.Command("iptables", "-P", "OUTPUT", "ACCEPT").CombinedOutput(); err != nil {
		return fmt.Errorf("defaulting OUTPUT to ACCEPT %v: %v", err, string(out))
	}

	return nil
}

func createEstablishedChain() error {
	if out, err := exec.Command("iptables", "-N", "ESTABLISHED").CombinedOutput(); err != nil {
		return fmt.Errorf("creating custom ESTABLISHED chain %v: %v", err, string(out))
	}

	rule := "-A ESTABLISHED -p %s -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT"

	for _, proto := range []string{"tcp", "udp", "icmp"} {
		args := strings.Fields(fmt.Sprintf(rule, proto))

		if out, err := exec.Command("iptables", args...).CombinedOutput(); err != nil {
			return fmt.Errorf("creating custom ESTABLISHED chain %v: %v", err, string(out))
		}
	}

	if out, err := exec.Command("iptables", "-A", "ESTABLISHED", "-j", "RETURN").CombinedOutput(); err != nil {
		return fmt.Errorf("creating custom ESTABLISHED chain %v: %v", err, string(out))
	}

	if out, err := exec.Command("iptables", "-A", "FORWARD", "-j", "ESTABLISHED").CombinedOutput(); err != nil {
		return fmt.Errorf("creating custom ESTABLISHED chain %v: %v", err, string(out))
	}

	return nil
}

func defaultForwardDrop() error {
	if out, err := exec.Command("iptables", "-P", "FORWARD", "DROP").CombinedOutput(); err != nil {
		return fmt.Errorf("defaulting FORWARD to DROP %v: %v", err, string(out))
	}

	return nil
}
