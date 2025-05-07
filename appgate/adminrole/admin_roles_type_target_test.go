package adminrole

import (
	"encoding/json"
	"testing"
)

var actionMapJSONResponse = `
{
	"All": [
		{
			"name": "All",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "AdministrativeRole",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdminMessage",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "AllocatedIp",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ApplianceCustomization",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AuditLog",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "AutoUpdate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Blacklist",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "CaCertificate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "DeviceClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Entitlement",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "FailedAuthentication",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Fido2Device",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "File",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "GlobalSetting",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "IpPool",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "License",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "LocalUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "OtpSeed",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Policy",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RegisteredDevice",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "RiskModel",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "RingfenceRule",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ServiceUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "SessionInfo",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Site",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "TokenRecord",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "TrustedCertificate",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "UserLicense",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "UserClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Ztp",
			"scopable": false,
			"scopableByIdp": false
		}
	],
	"AssignFunction": [
		{
			"name": "All",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": false,
			"scopableByIdp": false
		}
	],
	"Backup": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"CheckStatus": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"Create": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdministrativeRole",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ApplianceCustomization",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Blacklist",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "CaCertificate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "DeviceClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Entitlement",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "File",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "IpPool",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "LocalUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Policy",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RingfenceRule",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ServiceUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Site",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "TrustedCertificate",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "UserClaimScript",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"Delete": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdministrativeRole",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdminMessage",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ApplianceCustomization",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Blacklist",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "CaCertificate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "DeviceClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Entitlement",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Fido2Device",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "File",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "GlobalSetting",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "IpPool",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "License",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "LocalUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "OtpSeed",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Policy",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RegisteredDevice",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "RingfenceRule",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ServiceUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Site",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "TrustedCertificate",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "UserLicense",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "UserClaimScript",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"DownloadLogs": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"Edit": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdministrativeRole",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ApplianceCustomization",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AutoUpdate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "CaCertificate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "DeviceClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Entitlement",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "File",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "GlobalSetting",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "IpPool",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "License",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "LocalUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Policy",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RingfenceRule",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RiskModel",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ServiceUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Site",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "TrustedCertificate",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "UserClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Ztp",
			"scopable": false,
			"scopableByIdp": false
		}
	],
	"Export": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"GetUserAttributes": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"Reboot": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"RenewCertificate": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"Reevaluate": [
		{
			"name": "All",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "TokenRecord",
			"scopable": false,
			"scopableByIdp": true
		}
	],
	"Revoke": [
		{
			"name": "All",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "CaCertificate",
			"scopable": false,
			"scopableByIdp": false
		}
	],
	"Tag": [
		{
			"name": "All",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "AdministrativeRole",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ApplianceCustomization",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "DeviceClaimScript",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Entitlement",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IpPool",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "LocalUser",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Policy",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "RingfenceRule",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ServiceUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Site",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "TrustedCertificate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "UserClaimScript",
			"scopable": false,
			"scopableByIdp": false
		}
	],
	"Test": [
		{
			"name": "All",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Policy",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "UserClaimScript",
			"scopable": false,
			"scopableByIdp": false
		}
	],
	"Upgrade": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		}
	],
	"View": [
		{
			"name": "All",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdministrativeRole",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AdminMessage",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "AllocatedIp",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Appliance",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "ApplianceCustomization",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "AuditLog",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "AutoUpdate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "Blacklist",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "CaCertificate",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ClientProfile",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Condition",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "CriteriaScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "DeviceClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Entitlement",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "EntitlementScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "FailedAuthentication",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Fido2Device",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "File",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "GlobalSetting",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "IdentityProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "IpPool",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "License",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "LocalUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "MfaProvider",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "OtpSeed",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Policy",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RegisteredDevice",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "RingfenceRule",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "RiskModel",
			"scopable": false,
			"scopableByIdp": false
		},
		{
			"name": "ServiceUser",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "SessionInfo",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "Site",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "TokenRecord",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "TrustedCertificate",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "UserLicense",
			"scopable": false,
			"scopableByIdp": true
		},
		{
			"name": "UserClaimScript",
			"scopable": true,
			"scopableByIdp": false
		},
		{
			"name": "Ztp",
			"scopable": false,
			"scopableByIdp": false
		}
	]
}
`

func TestCanScopePrivlige(t *testing.T) {
	var actionMapJSON map[string]interface{}
	if err := json.Unmarshal([]byte(actionMapJSONResponse), &actionMapJSON); err != nil {
		t.Fatalf("could not resolve teststub %s", err)
	}
	type args struct {
		actionmap     map[string]interface{}
		privilegeType string
		target        string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Condition test",
			args: args{
				actionmap:     actionMapJSON,
				target:        "Condition",
				privilegeType: "Test",
			},
			want: false,
		},
		{
			name: "All tags",
			args: args{
				actionmap:     actionMapJSON,
				target:        "All",
				privilegeType: "Tag",
			},
			want: false,
		},
		{
			name: "All",
			args: args{
				actionmap:     actionMapJSON,
				target:        "All",
				privilegeType: "All",
			},
			want: false,
		},
		{
			name: "All AdministrativeRole",
			args: args{
				actionmap:     actionMapJSON,
				target:        "AdministrativeRole",
				privilegeType: "All",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanScopePrivlige(tt.args.actionmap, tt.args.privilegeType, tt.args.target); got != tt.want {
				t.Errorf("CanScopePrivlige() = %v, want %v", got, tt.want)
			}
		})
	}
}
