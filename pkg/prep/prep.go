package prep

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func importPrivateKeyFile() ssh.AuthMethod {
	buffer, err := ioutil.ReadFile("/etc/captain/builder/key.private")
	if err != nil {
		log.Println(err)
		log.Println("unable to read private key file")
		return nil
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Println(err)
		log.Println("unable to parse private key file")
		return nil
	}
	return ssh.PublicKeys(key)
}

func connectToPlane(hostname string) (ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			importPrivateKeyFile(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	connection, err := ssh.Dial("tcp", hostname + ":22", sshConfig)
	if err != nil {
		log.Println(err)
		return ssh.Session{}, errors.New("unable to dial SSH server")
	}
	session, err := connection.NewSession()
	if err != nil {
		log.Println(err)
		return ssh.Session{}, errors.New("unable to create remote SSH session")
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Println(err)
		session.Close()
		return ssh.Session{}, errors.New("request for pseudo terminal failed")
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Println(err)
		session.Close()
		return ssh.Session{}, errors.New("unable to setup stdin")
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Println(err)
		session.Close()
		return ssh.Session{}, errors.New("unable to setup stdout")
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := session.StderrPipe()
	if err != nil {
		log.Println(err)
		session.Close()
		return ssh.Session{}, errors.New("unable to setup stderr")
	}
	go io.Copy(os.Stderr, stderr)
	return *session, nil
}
// TODO: Fix PowerDNS so IP is no longer needed
func DeployPlan(hostname string, ip string, commands []string) error {
	for index, element := range commands {
		log.Println(fmt.Sprintf("%s [%d]: %s", hostname, index, element))
		session, err := connectToPlane(ip)
		if err != nil {
			return errors.New("unable to connect to remote plane")
		}
		err = session.Run(element)
		if err != nil {
			log.Println(err)
			session.Close()
			return errors.New("pre-flight prep failed")
		}
		session.Close()
	}
	return nil
}
