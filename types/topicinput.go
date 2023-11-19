package types

// `junjo` doesn't need to know the type you using
// `TopicInput`
type TopicInput interface {
	Validate(data map[string]interface{}) (bool, error)
}
