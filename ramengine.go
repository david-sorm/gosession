package gosession

type session map[string]interface{}

// A very simple implementation of GosessionEngine, doesn't support saving,
// hence the name "RamEngine"
type RamEngine struct {
	sessions map[string]session
	es EngineState
}

func (re RamEngine) GetEngineStatePointer() *EngineState {
	return &re.es
}

func(re RamEngine) Init() error {
	re.sessions = make(map[string]session)
	return nil
}

func (re RamEngine) Close() {
	re.sessions = nil
}

func (re RamEngine) SessionExists(sessionID string) bool {
	_, exists := re.sessions[sessionID]
	return exists
}

func (re RamEngine) CreateSession(sessionID string) {
	re.sessions[sessionID] = make(session)
}

func (re RamEngine) DestroySession(sessionID string) {
	delete(re.sessions, sessionID)
}

func (re RamEngine) DestroyAllSessions() {
	re.sessions = make(map[string]session)
}

func (re RamEngine) ReadKey(sessionID string, key string) interface{} {
	return re.sessions[sessionID][key]
}

func (re RamEngine) WriteKey(sessionID string, key string, Value interface{}) {
	re.sessions[sessionID][key] = Value
}

func (re RamEngine) DeleteKey(sessionID string, key string) {
	delete(re.sessions[sessionID], key)
}

