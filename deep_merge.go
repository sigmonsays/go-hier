package hier

import "fmt"

// deepMerge merges src into dst recursively.
func deepMerge(dst, src map[string]any) map[string]any {
	for k, srcVal := range src {
		if dstVal, exists := dst[k]; exists {
			dstMap, dstIsMap := dstVal.(map[string]any)
			srcMap, srcIsMap := srcVal.(map[string]any)

			if dstIsMap && srcIsMap {
				dst[k] = deepMerge(dstMap, srcMap)
				continue
			}
		}
		dst[k] = srcVal
	}
	return dst
}

// collectMerge walks an arbitrary JSON-shaped value (map, array, or scalar)
// and merges every map it finds into result. Arrays are walked element by
// element regardless of nesting depth; scalars are ignored.
func collectMerge(result map[string]any, node any) map[string]any {
	switch v := node.(type) {
	case map[string]any:
		result = deepMerge(result, v)
	case []any:
		for _, item := range v {
			result = collectMerge(result, item)
		}
	default:
		// scalar (string, number, bool, nil) at this level — nothing to merge
	}
	return result
}

// MergeNestedArraysAny handles input unmarshaled generically as `any`,
// regardless of how deeply/irregularly it's nested.
func MergeNestedArraysAny(input any) (map[string]any, error) {
	result := make(map[string]any)

	switch input.(type) {
	case []any, map[string]any:
		result = collectMerge(result, input)
	default:
		return nil, fmt.Errorf("expected array or object at top level, got %T", input)
	}

	return result, nil
}
