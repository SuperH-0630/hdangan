package v1main

type Closer interface {
	Close()
}

func WinClose(n Closer) {
	if n != nil {
		n.Close()
	}
}
