package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	inputJSON := `{
        "number_1": {
            "N": "1.50"
        },
        "string_1": {
            "S": "784498 "
        },
        "string_2": {
            "S": "2014-07-16T20:55:46Z"
        },
        "map_1": {
            "M": {
                "bool_1": {
                    "BOOL": "truthy"
                },
                "null_1": {
                    "NULL ": "true"
                },
                "list_1": {
                    "L": [
                        {
                            "S": ""
                        },
                        {
                            "N": "011"
                        },
                        {
                            "N": "5215s"
                        },
                        {
                            "BOOL": "f"
                        },
                        {
                            "NULL": "0"
                        }
                    ]
                }
            }
        },
        "list_2": {
            "L": "noop"
        },
        "list_3": {
            "L": [
                "noop"
            ]
        },
        "": {
            "S": "noop"
        }
    }`

	var input map[string]interface{}
	err := json.Unmarshal([]byte(inputJSON), &input)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	output := transformInput(input)
	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}

	fmt.Println(string(outputJSON))
}

func transformInput(input map[string]interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

	for key, value := range input {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		switch valueTyped := value.(type) {
		case map[string]interface{}:
			for dataType, val := range valueTyped {
				dataType = strings.TrimSpace(dataType)
				switch dataType {
				case "S":
					strVal := strings.TrimSpace(val.(string))
					if strVal == "" {
						continue
					}
					if t, err := time.Parse(time.RFC3339, strVal); err == nil {
						item[key] = t.Unix()
					} else {
						item[key] = strVal
					}
				case "N":
					if num, err := strconv.ParseFloat(strings.TrimSpace(val.(string)), 64); err == nil {
						item[key] = num
					}
				case "BOOL":
					boolVal, err := strconv.ParseBool(strings.TrimSpace(val.(string)))
					if err == nil {
						item[key] = boolVal
					}
				case "NULL":
					if boolVal, err := strconv.ParseBool(strings.TrimSpace(val.(string))); err == nil && boolVal {
						item[key] = nil
					}
				case "L":
					// Check if val is actually a slice before proceeding
					if listVal, ok := val.([]interface{}); ok {
						list := transformList(listVal)
						if len(list) > 0 {
							item[key] = list
						}
					}
				case "M":
					if mapVal, ok := val.(map[string]interface{}); ok {
						transformedMap := transformInput(mapVal)
						if len(transformedMap) > 0 {
							item[key] = transformedMap[0] // Assuming only one map item for simplicity
						}
					}
				}
			}
		}
	}

	if len(item) > 0 {
		result = append(result, item)
	}

	return result
}

func transformList(input []interface{}) []interface{} {
	list := make([]interface{}, 0)
	for _, element := range input {
		switch elemTyped := element.(type) {
		case map[string]interface{}:
			for dataType, val := range elemTyped {
				dataType = strings.TrimSpace(dataType)
				switch dataType {
				case "S":
					strVal := strings.TrimSpace(val.(string))
					if strVal != "" {
						list = append(list, strVal)
					}
				case "N":
					if num, err := strconv.ParseFloat(strings.TrimSpace(val.(string)), 64); err == nil {
						list = append(list, num)
					}
				case "BOOL":
					boolVal, err := strconv.ParseBool(strings.TrimSpace(val.(string)))
					if err == nil {
						list = append(list, boolVal)
					}
				}
			}
		}
	}
	return list
}
