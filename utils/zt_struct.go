package main

import (
	"fmt"
	"time"
)

type ZtNetPost struct {
	Description string `json:"description,omitempty"`
	Config      struct {
		Name            string `json:"name,omitempty"`
		Private         bool   `json:"private,omitempty"`
		EnableBroadcast bool   `json:"enableBroadcast,omitempty"`
		MulticastLimit  int    `json:"multicastLimit,omitempty"`
		Mtu             int    `json:"mtu,omitempty"`

		Routes []struct {
			Target string      `json:"target,omitempty"`
			Via    interface{} `json:"via,omitempty"`
		} `json:"routes,omitempty"`

		IPAssignmentPools []struct {
			IPRangeStart string `json:"ipRangeStart,omitempty"`
			IPRangeEnd   string `json:"ipRangeEnd,omitempty"`
		} `json:"ipAssignmentPools,omitempty"`

		V4AssignMode struct {
			Zt bool `json:"zt,omitempty"`
		} `json:"v4AssignMode,omitempty"`

		V6AssignMode struct {
			SixPlane bool `json:"6plane,omitempty"`
			Rfc4193  bool `json:"rfc4193,omitempty"`
			Zt       bool `json:"zt,omitempty"`
		} `json:"v6AssignMode,omitempty"`
	} `json:"config,omitempty"`
}

type ZtNetInfo struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`

	Config struct {
		Name            string `json:"name,omitempty"`
		Private         bool   `json:"private,omitempty"`
		EnableBroadcast bool   `json:"enableBroadcast,omitempty"`
		MulticastLimit  int    `json:"multicastLimit,omitempty"`
		Mtu             int    `json:"mtu,omitempty"`
		CreationTime    int64  `json:"creationTime,omitempty"`
		LastModified    int64  `json:"lastModified,omitempty"`

		DNS struct {
			Domain  string   `json:"domain,omitempty"`
			Servers []string `json:"servers,omitempty"`
		} `json:"dns,omitempty"`

		IPAssignmentPools []struct {
			IPRangeStart string `json:"ipRangeStart,omitempty"`
			IPRangeEnd   string `json:"ipRangeEnd,omitempty"`
		} `json:"ipAssignmentPools,omitempty"`

		Routes []struct {
			Target string      `json:"target,omitempty"`
			Via    interface{} `json:"via,omitempty"`
		} `json:"routes,omitempty"`

		V4AssignMode struct {
			Zt bool `json:"zt,omitempty"`
		} `json:"v4AssignMode,omitempty"`

		V6AssignMode struct {
			SixPlane bool `json:"6plane,omitempty"`
			Rfc4193  bool `json:"rfc4193,omitempty"`
			Zt       bool `json:"zt,omitempty"`
		} `json:"v6AssignMode,omitempty"`
	} `json:"config,omitempty"`

	OwnerID               string `json:"ownerId,omitempty"`
	OnlineMemberCount     int    `json:"onlineMemberCount,omitempty"`
	AuthorizedMemberCount int    `json:"authorizedMemberCount,omitempty"`
	TotalMemberCount      int    `json:"totalMemberCount,omitempty"`
}

type ZtNetMemberPost struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Hidden      bool   `json:"hidden,omitempty"`

	Config struct {
		Authorized      bool     `json:"authorized,omitempty"`
		ActiveBridge    bool     `json:"activeBridge,omitempty"`
		NoAutoAssignIps bool     `json:"noAutoAssignIps,omitempty"`
		IPAssignments   []string `json:"ipAssignments,omitempty"`
	} `json:"config,omitempty"`
}

type ZtNetMemberInfo struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	NetworkID       string `json:"networkId,omitempty"`
	NodeID          string `json:"nodeId,omitempty"`
	Hidden          bool   `json:"hidden,omitempty"`
	PhysicalAddress string `json:"physicalAddress,omitempty"`
	ClientVersion   string `json:"clientVersion,omitempty"`
	ProtocolVersion int    `json:"protocolVersion,omitempty"`
	Clock           int64  `json:"clock,omitempty"`
	LastOnline      int64  `json:"lastOnline,omitempty"`

	Config struct {
		ActiveBridge         bool     `json:"activeBridge,omitempty"`
		Authorized           bool     `json:"authorized,omitempty"`
		NoAutoAssignIps      bool     `json:"noAutoAssignIps,omitempty"`
		IPAssignments        []string `json:"ipAssignments,omitempty"`
		CreationTime         int64    `json:"creationTime,omitempty"`
		LastAuthorizedTime   int64    `json:"lastAuthorizedTime,omitempty"`
		LastDeauthorizedTime int      `json:"lastDeauthorizedTime,omitempty"`
	} `json:"config,omitempty"`
}

func (i ZtNetInfo) show_header() {
	fmt.Printf("%-18s %-16s %-10s %-10s %s\n", "NID", "Name", "Private", "O/T/A", "CreationTime")
}

func (i *ZtNetInfo) show() {
	createTime := time.Unix(i.Config.CreationTime/1000, 0).Local().Format(time.RFC3339)
	ota := fmt.Sprintf("%d/%d/%d", i.OnlineMemberCount, i.TotalMemberCount, i.AuthorizedMemberCount)
	fmt.Printf("%-18s %-16s %-10t %-10s %s\n", i.ID, i.Config.Name, i.Config.Private, ota, createTime)
}

func (i ZtNetMemberInfo) show_header() {
	fmt.Printf("%-12s %-14s %-18s %-18s %-10s %-12s %-7s %s\n",
		"MID", "Name", "IP_assign", "IP_physical", "Version", "LastOnline", "Auth", "Hidden")
}

func (i *ZtNetMemberInfo) show() {
	lastOnline := time.Unix(i.LastOnline/1000, 0)
	lastduration := time.Since(lastOnline).Truncate(time.Second)

	ip_assigned := "-"
	if len(i.Config.IPAssignments) > 0 {
		ip_assigned = i.Config.IPAssignments[0]
	}

	fmt.Printf("%-12s %-14s %-18s %-18s %-10s %-12s %-7t %t\n",
		i.NodeID, i.Name, ip_assigned, i.PhysicalAddress,
		i.ClientVersion, lastduration, i.Config.Authorized, i.Hidden)
}

func display_networks(networks []ZtNetInfo) {
	if len(networks) == 0 {
		return
	}

	if args.Verbose {
		for i, net := range networks {
			fmt.Printf("-- net %d: %s\n", i, net.ID)
			fmt.Println(dumps(net, args.Format))
		}
	} else {
		// show header
		ZtNetInfo{}.show_header()

		// show brief info
		for _, net := range networks {
			net.show()
		}
	}
}

func display_network_members(members []ZtNetMemberInfo) {
	if len(members) == 0 {
		return
	}

	if args.Verbose {
		for i, m := range members {
			fmt.Printf("-- netm %d: %s\n", i, m.NodeID)
			fmt.Println(dumps(m, args.Format))
		}
	} else {
		// show header
		ZtNetMemberInfo{}.show_header()

		// show brief info
		for _, m := range members {
			m.show()
		}
	}
}
