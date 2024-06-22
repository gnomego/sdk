package exec

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/gobuffalo/flect"
)

type SplatBool struct {
	TrueValue  string
	FalseValue string
}

type SplatOptions struct {
	Command         *[]string
	Prefix          *string
	Aliases         *map[string]string
	Assign          *string
	PreserveCase    *bool
	ShortFlag       *bool
	NonFlags        *[]string
	NonFlagValues   *SplatBool
	Includes        *[]string
	Excludes        *[]string
	Args            *[]string
	AppendArgs      *[]string
	ArrayDelimiter  *string
	ArrayConcatArgs *[]string
	PostArgs        *[]string
}

func SplatMap(m map[string]interface{}, options *SplatOptions) ([]string, error) {
	splat := []string{}

	prependArgs := []string{}
	appendArgs := []string{}
	postArgs := []string{}

	if options == nil {
		options = &SplatOptions{}
		*options.ShortFlag = true
		*options.Prefix = "-"
	}

	if options.Command != nil {
		splat = append(splat, *options.Command...)
	}

	for k, v := range m {
		name := k
		value := v
		t := reflect.TypeOf(v)
		kind := t.Kind()

		if options.Includes != nil {
			if slices.Contains(*options.Includes, k) {
				continue
			}
		}
		if options.Excludes != nil {
			if slices.Contains(*options.Excludes, k) {
				continue
			}
		}

		paramName := k
		if options.Aliases != nil && (*options.Aliases)[k] != "" {
			paramName = (*options.Aliases)[k]
		} else {
			if options.PreserveCase != nil && !*options.PreserveCase {
				paramName = flect.Dasherize(flect.Underscore(k))
			}
			if options.ShortFlag != nil && *options.ShortFlag {
				paramName = *options.Prefix + paramName
			}
		}

		if kind == reflect.Array || kind == reflect.Slice {
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, value.([]string)...)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				postArgs = append(postArgs, value.([]string)...)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, value.([]string)...)
				continue
			}

			if options.ArrayConcatArgs != nil && slices.Contains(*options.ArrayConcatArgs, name) {
				if options.ArrayDelimiter == nil {
					*options.ArrayDelimiter = ","
				}

				if options.Assign != nil {
					v := paramName + *options.Assign + "\"" + strings.Join(value.([]string), *options.ArrayDelimiter) + "\""
					splat = append(splat, v)
				} else {
					splat = append(splat, paramName, strings.Join(value.([]string), *options.ArrayDelimiter))
				}

				continue
			}

			for _, v := range value.([]string) {
				splat = append(splat, paramName, v)
			}
		}

		if kind == reflect.String {
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, value.(string))
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, value.(string))
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, value.(string))
				continue
			}
		}

		// if value is bool
		if kind == reflect.Bool {
			if options.NonFlags != nil && slices.Contains(*options.NonFlags, name) {
				sp := options.NonFlagValues
				if sp == nil {
					sp = &SplatBool{
						TrueValue:  "true",
						FalseValue: "false",
					}
				}

				if value == true {
					if options.Assign != nil {
						v := paramName + *options.Assign + sp.TrueValue
						splat = append(splat, v)
					} else {
						splat = append(splat, paramName, sp.TrueValue)
					}
					continue
				} else {
					if options.Assign != nil {
						v := paramName + *options.Assign + sp.FalseValue
						splat = append(splat, v)
					} else {
						splat = append(splat, paramName, sp.FalseValue)
					}
					continue
				}
			}

			if value == true {
				splat = append(splat, paramName)
			}

			continue
		}

		if kind == reflect.Int {
			iv := value.(int)
			v1 := fmt.Sprintf("%d", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		if kind == reflect.Float64 {
			iv := value.(float64)
			v1 := fmt.Sprintf("%f", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		if kind == reflect.Int16 {
			iv := value.(int16)
			v1 := fmt.Sprintf("%d", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		if kind == reflect.Int32 {
			iv := value.(int32)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, string(iv))
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, string(iv))
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, string(iv))
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + string(iv)
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, string(iv))
			}
		}

		if kind == reflect.Int64 {
			iv := value.(int64)
			v1 := fmt.Sprintf("%d", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		return []string{}, fmt.Errorf("unsupported type %v", kind)
	}

	if options.Command != nil && len(*options.Command) > 0 {
		splat = append(*options.Command, splat...)
	}

	if prependArgs != nil && len(prependArgs) > 0 {
		splat = append(prependArgs, splat...)
	}

	if appendArgs != nil && len(appendArgs) > 0 {
		splat = append(splat, appendArgs...)
	}

	if postArgs != nil && len(postArgs) > 0 {
		splat = append(splat, "--")
		splat = append(splat, postArgs...)
	}

	return splat, nil
}

func Splat(object interface{}, options *SplatOptions) ([]string, error) {
	t := reflect.TypeOf(object)
	if t.Kind() == reflect.Map {
		return SplatMap(object.(map[string]interface{}), options)
	}

	prependArgs := []string{}
	appendArgs := []string{}
	postArgs := []string{}

	splat := []string{}

	if options == nil {
		options = &SplatOptions{}
		*options.ShortFlag = true
		*options.Prefix = "-"
	}

	if options.Command != nil {
		splat = append(splat, *options.Command...)
	}

	// loop through the object properties
	// and add them to the splat
	// if they are not in the excludes list
	// and are in the includes list
	// if the includes list is not empty
	// if the object is a map
	// if the object is a struct
	// if the object is a pointer
	v := reflect.ValueOf(object)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		paramName := name
		value := v.Field(i).Interface()
		kind := field.Type.Kind()
		if options.Includes != nil {
			if slices.Contains(*options.Includes, name) {
				continue
			}
		}
		if options.Excludes != nil {
			if slices.Contains(*options.Excludes, name) {
				continue
			}
		}

		if options.Aliases != nil && (*options.Aliases)[name] != "" {
			paramName = (*options.Aliases)[name]
		} else {
			if options.PreserveCase != nil && !*options.PreserveCase {
				paramName = flect.Dasherize(flect.Underscore(name))
			}
			if options.ShortFlag != nil && *options.ShortFlag {
				paramName = *options.Prefix + paramName
			}
		}

		if kind == reflect.Array || kind == reflect.Slice {
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, value.([]string)...)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				postArgs = append(postArgs, value.([]string)...)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, value.([]string)...)
				continue
			}

			if options.ArrayConcatArgs != nil && slices.Contains(*options.ArrayConcatArgs, name) {
				if options.ArrayDelimiter == nil {
					*options.ArrayDelimiter = ","
				}

				if options.Assign != nil {
					v := paramName + *options.Assign + "\"" + strings.Join(value.([]string), *options.ArrayDelimiter) + "\""
					splat = append(splat, v)
				} else {
					splat = append(splat, paramName, strings.Join(value.([]string), *options.ArrayDelimiter))
				}

				continue
			}

			for _, v := range value.([]string) {
				splat = append(splat, paramName, v)
			}
		}

		if kind == reflect.String {
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, value.(string))
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, value.(string))
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, value.(string))
				continue
			}
		}

		// if value is bool
		if kind == reflect.Bool {
			if options.NonFlags != nil && slices.Contains(*options.NonFlags, name) {
				sp := options.NonFlagValues
				if sp == nil {
					sp = &SplatBool{
						TrueValue:  "true",
						FalseValue: "false",
					}
				}

				if value == true {
					if options.Assign != nil {
						v := paramName + *options.Assign + sp.TrueValue
						splat = append(splat, v)
					} else {
						splat = append(splat, paramName, sp.TrueValue)
					}
					continue
				} else {
					if options.Assign != nil {
						v := paramName + *options.Assign + sp.FalseValue
						splat = append(splat, v)
					} else {
						splat = append(splat, paramName, sp.FalseValue)
					}
					continue
				}
			}

			if value == true {
				splat = append(splat, paramName)
			}

			continue
		}

		if kind == reflect.Int {
			iv := value.(int)
			v1 := fmt.Sprintf("%d", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		if kind == reflect.Float64 {
			iv := value.(float64)
			v1 := fmt.Sprintf("%f", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		if kind == reflect.Int16 {
			iv := value.(int16)
			v1 := fmt.Sprintf("%d", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		if kind == reflect.Int32 {
			iv := value.(int32)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, string(iv))
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, string(iv))
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, string(iv))
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + string(iv)
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, string(iv))
			}
		}

		if kind == reflect.Int64 {
			iv := value.(int64)
			v1 := fmt.Sprintf("%d", iv)
			if options.Args != nil && slices.Contains(*options.Args, name) {
				prependArgs = append(prependArgs, v1)
				continue
			}

			if options.AppendArgs != nil && slices.Contains(*options.AppendArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.PostArgs != nil && slices.Contains(*options.PostArgs, name) {
				appendArgs = append(appendArgs, v1)
				continue
			}

			if options.Assign != nil {
				v := paramName + *options.Assign + v1
				splat = append(splat, v)
			} else {
				splat = append(splat, paramName, v1)
			}
		}

		return []string{}, fmt.Errorf("unsupported type %v", kind)
	}

	if options.Command != nil && len(*options.Command) > 0 {
		splat = append(*options.Command, splat...)
	}

	if prependArgs != nil && len(prependArgs) > 0 {
		splat = append(prependArgs, splat...)
	}

	if appendArgs != nil && len(appendArgs) > 0 {
		splat = append(splat, appendArgs...)
	}

	if postArgs != nil && len(postArgs) > 0 {
		splat = append(splat, "--")
		splat = append(splat, postArgs...)
	}

	return splat, nil
}
