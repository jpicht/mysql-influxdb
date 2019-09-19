package app

type (
	TmpPoint struct {
		tags   map[string]string
		values map[string]interface{}
	}
)

func newTmpPointProcList(tags map[string]string, user string, command string, db string, host string, state string) *TmpPoint {
	copy := make(map[string]string)
	for k, v := range tags {
		copy[k] = v
	}
	copy["user"] = user
	copy["command"] = command
	copy["db"] = db
	copy["client"] = host
	copy["state"] = state
	return &TmpPoint{
		tags:   copy,
		values: make(map[string]interface{}),
	}
}
