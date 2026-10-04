package ipguard

type LocalSocket struct {
	out chan<- [][]byte
	in  <-chan [][]byte
}

func createPipe() (*LocalSocket, *LocalSocket) {
	ch1 := make(chan [][]byte, 5)
	ch2 := make(chan [][]byte, 5)
	s1 := createSocket(ch1, ch2)
	s2 := createSocket(ch2, ch1)
	return s1, s2
}

func createSocket(in <-chan [][]byte, out chan<- [][]byte) *LocalSocket {
	socket := LocalSocket{}
	socket.in = in
	socket.out = out
	return &socket
}

func (s *LocalSocket) Bind(string) (int, error) {
	return 0, nil
}

func (s *LocalSocket) SendMessage(parts [][]byte) error {
	s.out <- parts
	return nil
}
func (s *LocalSocket) RecvMessage() ([][]byte, error) {
	return <-s.in, nil
}
func (s *LocalSocket) Identity() string {
	return ""
}
func (s *LocalSocket) SetIdentity(val string) {

}
func (s *LocalSocket) Connect(endpoint string) error {
	return nil
}
func (s *LocalSocket) SetCurveServerkey(val string) {

}
func (s *LocalSocket) SetCurveServer(int) {

}
