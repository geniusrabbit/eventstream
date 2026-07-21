package sql

// QueryConfig extra information
type QueryConfig struct {
	FieldObject any
	Target      string
	IterateBy   string
}

// QueryOption func
type QueryOption func(*QueryConfig)

// QWithMessageTmpl source info
func QWithMessageTmpl(fields any) QueryOption {
	return func(qc *QueryConfig) {
		qc.FieldObject = fields
	}
}

// QWithMessageTmplOf is a generic helper equivalent to QWithMessageTmpl(new(T)).
func QWithMessageTmplOf[T any]() QueryOption {
	return QWithMessageTmpl(new(T))
}

// QWithTarget data collection
func QWithTarget(target string) QueryOption {
	return func(qc *QueryConfig) {
		qc.Target = target
	}
}

// QWithIterateBy message value
func QWithIterateBy(iterateBy string) QueryOption {
	return func(qc *QueryConfig) {
		qc.IterateBy = iterateBy
	}
}
