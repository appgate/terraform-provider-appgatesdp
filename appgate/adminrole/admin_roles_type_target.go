package adminrole

type RoleScope struct {
	Name          string `json:"name"`
	Scopable      bool   `json:"scopable"`
	ScopableByIdp bool   `json:"scopableByIdp"`
}

func (r *RoleScope) CanUseScope() bool {
	return r.Scopable || r.ScopableByIdp
}
func toRoleScopeMap(v map[string]interface{}) map[string][]RoleScope {
	roleScopes := make(map[string][]RoleScope)
	for key, value := range v {
		if valueMap, ok := value.([]interface{}); ok {
			for _, item := range valueMap {
				if itemMap, ok := item.(map[string]interface{}); ok {
					rs := RoleScope{}
					if v, ok := itemMap["name"].(string); ok {
						rs.Name = v
					}
					if v, ok := itemMap["scopable"].(bool); ok {
						rs.Scopable = v
					}
					if v, ok := itemMap["scopableByIdp"].(bool); ok {
						rs.ScopableByIdp = v
					}
					roleScopes[key] = append(roleScopes[key], rs)
				}
			}
		}
	}
	return roleScopes
}

func CanScopePrivlige(v map[string]interface{}, privilegeType, target string) bool {
	actionmap := toRoleScopeMap(v)
	if list, ok := actionmap[privilegeType]; ok {
		for _, item := range list {
			if target == item.Name {
				return item.CanUseScope()
			}
		}
	}
	return false
}
