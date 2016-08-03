package main


import (
        "fmt"
        "net/url"
        "os"
        "io/ioutil"
        "log"
        "github.com/huin/goupnp"
        "github.com/huin/goupnp/dcps/internetgateway1"
        "github.com/huin/goupnp/dcps/internetgateway2"
)




type user_input interface {
        action() bool
}

type api_call struct {
        call_name string
        endpoint string
}

type local_command struct{
        cmd_name string
}


func (l local_command) action() bool {
    files, err := ioutil.ReadDir(".")
    if err != nil {
                        log.Fatal(err)
                                }
    fmt.Println(files)
        return true
}

func (a api_call) action() bool {
        clients, errors, err := internetgateway1.NewWANPPPConnection1Clients()
        extIPClients := make([]GetExternalIPAddresser, len(clients))
        for i, client := range clients {
                extIPClients[i] = client
        }
        DisplayExternalIPResults(extIPClients, errors, err)
        return true 

}

// Use discovered WANPPPConnection1 services to find external IP addresses.
func Example_WANPPPConnection1_GetExternalIPAddress() {
        clients, errors, err := internetgateway1.NewWANPPPConnection1Clients()
        extIPClients := make([]GetExternalIPAddresser, len(clients))
        for i, client := range clients {
                extIPClients[i] = client
        }
        DisplayExternalIPResults(extIPClients, errors, err)
        // Output:
}

// Use discovered WANIPConnection services to find external IP addresses.
func Example_WANIPConnection_GetExternalIPAddress() {
        clients, errors, err := internetgateway1.NewWANIPConnection1Clients()
        extIPClients := make([]GetExternalIPAddresser, len(clients))
        for i, client := range clients {
                extIPClients[i] = client
        }
        DisplayExternalIPResults(extIPClients, errors, err)
        // Output:
}

type GetExternalIPAddresser interface {
        GetExternalIPAddress() (NewExternalIPAddress string, err error)
        GetServiceClient() *goupnp.ServiceClient
}

func DisplayExternalIPResults(clients []GetExternalIPAddresser, errors []error, err error) {
        if err != nil {
                fmt.Fprintln(os.Stderr, "Error discovering service with UPnP: ", err)
                return
        }

        if len(errors) > 0 {
                fmt.Fprintf(os.Stderr, "Error discovering %d services:\n", len(errors))
                for _, err := range errors {
                        fmt.Println("  ", err)
                }
        }

        fmt.Fprintf(os.Stderr, "Successfully discovered %d services:\n", len(clients))
        for _, client := range clients {
                device := &client.GetServiceClient().RootDevice.Device

                fmt.Fprintln(os.Stderr, "  Device:", device.FriendlyName)
                if addr, err := client.GetExternalIPAddress(); err != nil {
                        fmt.Fprintf(os.Stderr, "    Failed to get external IP address: %v\n", err)
                } else {
                        fmt.Fprintf(os.Stderr, "    External IP address: %v\n", addr)
                }
        }
}

func Example_ReuseDiscoveredDevice() {
        var allMaybeRootDevices []goupnp.MaybeRootDevice
        for _, urn := range []string{internetgateway1.URN_WANPPPConnection_1, internetgateway1.URN_WANIPConnection_1} {
                maybeRootDevices, err := goupnp.DiscoverDevices(internetgateway1.URN_WANPPPConnection_1)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "Could not discover %s devices: %v\n", urn, err)
                }
                allMaybeRootDevices = append(allMaybeRootDevices, maybeRootDevices...)
        }
        locations := make([]*url.URL, 0, len(allMaybeRootDevices))
        fmt.Fprintf(os.Stderr, "Found %d devices:\n", len(allMaybeRootDevices))
        for _, maybeRootDevice := range allMaybeRootDevices {
                if maybeRootDevice.Err != nil {
                        fmt.Fprintln(os.Stderr, "  Failed to probe device at ", maybeRootDevice.Location.String())
                } else {
                        locations = append(locations, maybeRootDevice.Location)
                        fmt.Fprintln(os.Stderr, "  Successfully probed device at ", maybeRootDevice.Location.String())
                }
        }
        fmt.Fprintf(os.Stderr, "Attempt to re-acquire %d devices:\n", len(locations))
        for _, location := range locations {
                if _, err := goupnp.DeviceByURL(location); err != nil {
                        fmt.Fprintf(os.Stderr, "  Failed to reacquire device at %s: %v\n", location.String(), err)
                } else {
                        fmt.Fprintf(os.Stderr, "  Successfully reacquired device at %s\n", location.String())
                }
        }
        // Output:
}

// Use discovered igd1.WANCommonInterfaceConfig1 services to discover byte
// transfer counts.
func Example_WANCommonInterfaceConfig1_GetBytesTransferred() {
        clients, errors, err := internetgateway1.NewWANCommonInterfaceConfig1Clients()
        if err != nil {
                fmt.Fprintln(os.Stderr, "Error discovering service with UPnP:", err)
                return
        }
        fmt.Fprintf(os.Stderr, "Error discovering %d services:\n", len(errors))
        for _, err := range errors {
                fmt.Println("  ", err)
        }
        for _, client := range clients {
                if recv, err := client.GetTotalBytesReceived(); err != nil {
                        fmt.Fprintln(os.Stderr, "Error requesting bytes received:", err)
                } else {
                        fmt.Fprintln(os.Stderr, "Bytes received:", recv)
                }
                if sent, err := client.GetTotalBytesSent(); err != nil {
                        fmt.Fprintln(os.Stderr, "Error requesting bytes sent:", err)
                } else {
                        fmt.Fprintln(os.Stderr, "Bytes sent:", sent)
                }
        }
        // Output:
}

// Use discovered igd2.WANCommonInterfaceConfig1 services to discover byte
// transfer counts.
func Example_WANCommonInterfaceConfig2_GetBytesTransferred() {
        clients, errors, err := internetgateway2.NewWANCommonInterfaceConfig1Clients()
        if err != nil {
                fmt.Fprintln(os.Stderr, "Error discovering service with UPnP:", err)
                return
        }
        fmt.Fprintf(os.Stderr, "Error discovering %d services:\n", len(errors))
        for _, err := range errors {
                fmt.Println("  ", err)
        }
        for _, client := range clients {
                if recv, err := client.GetTotalBytesReceived(); err != nil {
                        fmt.Fprintln(os.Stderr, "Error requesting bytes received:", err)
                } else {
                        fmt.Fprintln(os.Stderr, "Bytes received:", recv)
                }
                if sent, err := client.GetTotalBytesSent(); err != nil {
                        fmt.Fprintln(os.Stderr, "Error requesting bytes sent:", err)
                } else {
                        fmt.Fprintln(os.Stderr, "Bytes sent:", sent)
                }
        }
        // Output:
}

func exec_cmd(u user_input) {
        u.action()
}



func main(){


    ls := local_command{cmd_name: "ls"}
    ip := api_call{endpoint: "", call_name: "external_ip"}

exec_cmd(ls)
exec_cmd(ip)


//turd := 
//Example_WANIPConnection_GetExternalIPAddress()
//fmt.Println(turd)

}
