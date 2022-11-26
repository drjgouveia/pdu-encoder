package main

import (
	"flag"
	"fmt"

	"gopkg.in/webnice/pdu.v1"
)

func main() {
	var err error
	var pduCoder pdu.Interface
	var x int
	sca := flag.String("sca", "+35193121314", "The SMS Gateway number")
	address := flag.String("address", "", "The SMS destination number")
	message := flag.String("message", "", "The message body to be sent on the SMS")
	flag.Parse()
	fmt.Printf("SCA: \"%s\"\nAddress: \"%s\"\nMessage: \"%s\"\n", *sca, *address, *message)
	var sms = pdu.Encode{
		Sca:                 *sca,
		Ucs2:                false,
		Flash:               false,
		Address:             *address,
		Message:             *message,
		StatusReportRequest: false,
	}
	pduCoder = pdu.New().Decoder(messageReceiver)
	defer pduCoder.Done()

	// Encode SMS
	var enc []string
	enc, err = pduCoder.Encoder(sms)
	if err != nil {
		println("Error encode message: %s", err.Error())
	}
	for x = range enc {
		println(enc[x])
	}
}

// Receive new messages
func messageReceiver(msg pdu.Message) {
	var out string

	out += "New message found\n"
	if msg.Error() != nil {
		out += "Message error\n"
		out += msg.Error().Error() + "\n"
		return
	}

	out += " SMSC:"
	out += msg.ServiceCentreAddress()
	if msg.Type() == pdu.TypeSmsStatusReport {
		out += " (status report: ["
		out += msg.DischargeTime().String()
		out += "] '"
		out += msg.ReportStatus().String()
		out += "')"
	}
	out += "\n"

	out += " From:"
	out += msg.OriginatingAddress() + "\n"
	out += " SMS ("
	out += fmt.Sprintf("%d", msg.DataParts())
	out += "):"
	out += msg.Data() + "\n"

	println(out)
}
