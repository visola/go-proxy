package upstream

// arrayOrString can be used to unmarshal values as string or array of strings
type arrayOrString []string

// UnmarshalYAML implement the unmarshal from YAML package
func (v *arrayOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		// Apparently value is not an array
		var single string
		err := unmarshal(&single)
		if err != nil {
			// Still can't parse it
			return err
		}
		*v = []string{single}
	} else {
		*v = multi
	}
	return nil
}

// FromMapOfArrayOfStrings converts a map of arrayOrString to a map of array of strings
func FromMapOfArrayOfStrings(original map[string]arrayOrString) map[string][]string {
	result := make(map[string][]string)
	for key, values := range original {
		result[key] = values
	}
	return result
}
