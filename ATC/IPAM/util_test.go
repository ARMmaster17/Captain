package IPAM

import (
	"net"
	"reflect"
	"testing"
)
/* There is a bug with how reflect.DeepEqual works with net.IPNet
https://stackoverflow.com/questions/42921227/ip-part-of-net-ipnet-does-not-fulfil-reflect-deepequal-but-are-equal
func Test_parseSubnetBlocks(t *testing.T) {
	type args struct {
		blocks []string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 []net.IPNet
	}{
		{
			name: "parses single address in /8 block",
			args: func(t *testing.T) args {
				return args{blocks: []string{"10.0.0.0/8"}}
			},
			want1: []net.IPNet{net.IPNet{
				IP: net.ParseIP("10.0.0.0"),
				Mask: net.IPMask{255,0,0,0},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := parseSubnetBlocks(tArgs.blocks)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseSubnetBlocks got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}*/

func Test_nextIP(t *testing.T) {
	type args struct {
		ip  net.IP
		inc uint
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 net.IP
	}{
		{
			name: "Test base /8",
			args: func (t *testing.T) args {
				return args{
					ip: net.ParseIP("10.0.0.1"),
					inc: 1,
				}
			},
			want1: net.ParseIP("10.0.0.2"),
		},
		{
			name: "Test base /8 rollover",
			args: func (t *testing.T) args {
				return args{
					ip: net.ParseIP("10.0.0.255"),
					inc: 1,
				}
			},
			want1: net.ParseIP("10.0.1.0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := nextIP(tArgs.ip, tArgs.inc)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("nextIP got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_subnetIsFull(t *testing.T) {
	type args struct {
		existingAddresses int64
		subnetCIDR        net.IPMask
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 bool
	}{
		{
			name: "test8+",
			args: func(t *testing.T) args {
				_, subnet, _ := net.ParseCIDR("10.0.0.0/8")
				return args{
					existingAddresses: 16777217,
					subnetCIDR: subnet.Mask,
				}
			},
			want1: true,
		},
		{
			name: "test8",
			args: func(t *testing.T) args {
				_, subnet, _ := net.ParseCIDR("10.0.0.0/8")
				return args{
					existingAddresses: 16777216,
					subnetCIDR: subnet.Mask,
				}
			},
			want1: true,
		},
		{
			name: "test8-",
			args: func(t *testing.T) args {
				_, subnet, _ := net.ParseCIDR("10.0.0.0/8")
				return args{
					existingAddresses: 16777215,
					subnetCIDR: subnet.Mask,
				}
			},
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := subnetIsFull(tArgs.existingAddresses, tArgs.subnetCIDR)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("subnetIsFull got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}

func Test_getSubnetAddressSize(t *testing.T) {
	type args struct {
		subnetCIDR net.IPMask
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 int64
	}{
		{
			name: "test30",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/30")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 4,
		},
		{
			name: "test29",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/29")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 8,
		},
		{
			name: "test28",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/28")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 16,
		},
		{
			name: "test27",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/27")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 32,
		},
		{
			name: "test26",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/26")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 64,
		},
		{
			name: "test25",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/25")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 128,
		},
		{
			name: "test24",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/24")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 256,
		},
		{
			name: "test23",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/23")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 512,
		},
		{
			name: "test22",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/22")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 1024,
		},
		{
			name: "test21",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/21")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 2048,
		},
		{
			name: "test20",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/20")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 4096,
		},
		// This is probably enough tests to get the point across.
		{
			name: "test8",
			args: func(t *testing.T) args {
				_, ipNet, _ := net.ParseCIDR("10.0.0.0/8")
				return args{subnetCIDR: ipNet.Mask}
			},
			want1: 16777216,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := getSubnetAddressSize(tArgs.subnetCIDR)

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getSubnetAddressSize got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
