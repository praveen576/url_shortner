package url_mapper_helpers

import(
)

func ValidateOpType(op_type string) bool {
	allowed_ops := []string { "shorten", "extend" }
	ok := false
	for _,val := range(allowed_ops) {
	  if op_type == val {
	    ok = true
	    break
	  }
	}
	return ok
}
