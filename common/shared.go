package common

// Run context
//
// Used in allmost all function as a first parameter
type DCtx struct {
	// Thread identifier
	ThreadId string
	// Job identifier
	JobId string
	// Environment name if the job refers to environment
	EnvName string
	// Service name if the job refers to service
	SrvName string
	// Project name if the job refers to project
	ProjName string
}

// Command result
type CmdRes struct {
	// Standard outtput content
	StdOut string
	// Standard error content
	StdErr string
	// Error which occurred
	Err error
}

// Multi command result
type CmdResMulti struct {
	// Error which occurred in some command
	Err error
	// Command result which produced an error
	ErrRes CmdRes
	// All command results
	Results []CmdRes
}

// Appends command result
func (m *CmdResMulti) Append(r CmdRes) {
	m.Results = append(m.Results, r)
	if r.Err != nil {
		m.Err = r.Err
		m.ErrRes = r
	}
}

// Append multi command result
func (m *CmdResMulti) AppendMulti(a CmdResMulti) {
	m.Results = append(m.Results, a.Results...)
	if a.Err != nil {
		m.Err = a.Err
		m.ErrRes = a.ErrRes
	}
}
