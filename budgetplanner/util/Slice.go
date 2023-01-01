package util

// AddToSlice adds values to their respective slices,
// eg:- columnName will be added to columnNames
func AddToSlice(columnName string, condition string, operator string, value interface{},
	columnNames *[]string, conditions *[]string, operators *[]string, values *[]interface{}) {
	// Operators should be one less than the rest.
	if len(*columnNames) != 0 {
		*operators = append(*operators, operator)
	}
	*columnNames = append(*columnNames, columnName)
	*conditions = append(*conditions, condition)
	*values = append(*values, value)
}
