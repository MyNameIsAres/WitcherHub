package main

import (
	"encoding/binary"
    "fmt"
    "net"
	"net/http"
	"time"
	"io/ioutil"
	"strings"
	"encoding/json"
	"math/rand"
)

type Const struct {
	TypeByte       []byte
	TypeStringUtf8 []byte
	TypeStringUtf16 []byte
	TypeInt64      []byte
	TypeInt16      []byte
	TypeUint32     []byte
	TypeInt32      []byte
	PacketHead     []byte
	PacketTail     []byte
	NsRemote 	   string
}

var constant = Const{
	TypeByte:       []byte{0x81, 0x08},
	TypeStringUtf8: []byte{0xAC, 0x08},
	TypeStringUtf16:[]byte{0x9C, 0x16},
	TypeInt64:       []byte{0x81, 0x64},
	TypeInt16:      []byte{0x81, 0x16},
	TypeUint32:     []byte{0x71, 0x32},
	TypeInt32:       []byte{0x81, 0x32},
	PacketHead:     []byte{0xDE, 0xAD},
	PacketTail:     []byte{0xBE, 0xEF},
	NsRemote:       "Remote",
}

type MapPinCoords struct {
	LocationX float64 `json:"locationx"`
	LocationY float64 `json:"locationy"`
}

type CompanionRequestBody struct {
	Username string `json:"username"`
	CompanionName string `json:"companionname"`
}

type BlockRequestBody struct {
	BlockedAction string `json:"blockedactionname"`
}

type GameData struct {
	IsActive string
}

type TwitchName struct {
	DisplayName string `json:"display_name"`
}

func Execute(command string) []byte {
	x := 0x81160008
	y := int32(x)

	payload := Init()

	payload = AppendUtf8(payload, constant.NsRemote)
	fmt.Println(payload)

	payload = AppendInt32(payload, int32(0x12345678))
	fmt.Println(payload)


	payload = AppendInt32(payload, int32(y))
	fmt.Println(payload)


	payload = AppendUtf8(payload, command)
	

	return End(payload)
}

func Init() []byte {
	return []byte{}
}

func End(payload []byte) []byte {
	lengthBytes := int16ToBytes(int16(len(payload) + 6))
	fmt.Println("LengthBytes", lengthBytes)
	return Append(Append(constant.PacketHead, lengthBytes), Append(payload, constant.PacketTail))
}


func AppendUtf8(payload []byte, text string) []byte {
	textLength := int16(len(text))
	foobar := AppendInt16([]byte{}, textLength)

	return Append(Append(Append(payload, constant.TypeStringUtf8), AppendInt16([]byte{}, textLength)), []byte(text))
}

func Append(payload []byte, data []byte) []byte {
	return append(payload, data...)
}
func AppendUtf16(payload []byte, text string) []byte {
	return Append(Append(Append(payload, constant.TypeUint32), int16ToBytes(int16(len(text)))), []byte(text))
}

func AppendByte(payload []byte, value byte) []byte {
	return Append(Append(payload, constant.TypeByte), []byte{value})
}

func int16ToBytes(value int16) []byte {
    bytes := make([]byte, 2)
    binary.BigEndian.PutUint16(bytes, uint16(value))

    return bytes
}

func int32ToBytes(value int32) []byte {
    bytes := make([]byte, 4)
    binary.BigEndian.PutUint32(bytes, uint32(value))

    return bytes
}

func AppendInt32(payload []byte, value int32) []byte {
	return Append(Append(payload, constant.TypeInt32), int32ToBytes(value))
}

func AppendInt16(payload []byte, value int16) []byte {
	return Append(Append(payload, constant.TypeInt16), int16ToBytes(value))
}


func getDisplayName(username string) (string, error) {
	clientID := ""
	bearerToken := ""

	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return username, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", bearerToken)

	q := req.URL.Query()
	q.Add("login", username)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return username, fmt.Errorf("error sending request: %w", err)
	}

	

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return username, fmt.Errorf("error reading response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return username, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}
	
	defer resp.Body.Close()

	var result struct {
		Data []TwitchName `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return username, fmt.Errorf("error parsing JSON response: %w", err)
	}
	fmt.Println("Result: ")
	fmt.Println(result)

	if len(result.Data) > 0 {
		displayName := result.Data[0].DisplayName
		fmt.Println(displayName) 
		return displayName, nil
	}

	return username, nil
}

// TODO Clean this mess up
func spawnMonsters(w http.ResponseWriter, r *http.Request) {
	// command := "spawn('Nekker', 3)"
	fmt.Println("Request came through");
	time.Sleep(5 * time.Second) 

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Unable to read request body", http.StatusBadRequest)
	// 	return
	// }

	// var companionReqBody CompanionRequestBody
	// err = json.Unmarshal(body, &companionReqBody)
	// 	if err != nil {
	// 	http.Error(w, "Invalid JSON format", http.StatusBadRequest)
	// 	return
	// }

		command := "spawn('Nekker', 3)"
	// command := "TwitchQuest(\"Missing Beast\", \"Find Gillys mother\", \"Amzie\")"


	// preCommand := "SpawnSpecialCompanion(\"Triss\", \"test\")"
	// command := strings.Replace(preCommand, "test", companionReqBody.Username, 1)
	fmt.Println(command)
	result := Execute(command)
		// Attempt to connect to the specified IP and port
		conn, err := net.Dial("tcp", "127.0.0.1:37001")
		if err != nil {
			fmt.Println("Connection failed:", err)
			return
		}
		defer conn.Close()
	
		fmt.Println("Connection successful!")

		// conn.

		
		
		// _, myErr := conn.Write(result)
		// if myErr != nil {
		// 	fmt.Println("Error writing:", myErr)
		// 	return
		// }
	
		fmt.Println("Data sent successfully")

		fmt.Print("Data sent successfully: [")
		for i, b := range result {
			fmt.Print(b)
			if i != len(result)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println("]")

}

func placeMapPin(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body) 
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var mapPinCoords MapPinCoords
	err = json.Unmarshal(body, &mapPinCoords) 
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	command := fmt.Sprintf("PlaceMapPin(\"%f\", \"%f\")", mapPinCoords.LocationX, mapPinCoords.LocationY)

	fmt.Println("Command: " + command);

	time.Sleep(3 * time.Second) 


	result := Execute(command)


	conn, err := net.Dial("tcp", "127.0.0.1:37001")
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connection successful!")
	_, myErr := conn.Write(result)
		if myErr != nil {
			fmt.Println("Error writing:", myErr)
			return
		}
	
		fmt.Println("[WitcherConnect] We placed a marker in the world")
}

func nakedGeralt(w http.ResponseWriter, r *http.Request) {
	command := "nakedGeralt()"
	fmt.Println("[WitcherConnect] A request to make Geralt naked came through!")
	time.Sleep(3 * time.Second) 

	result := Execute(command)
		conn, err := net.Dial("tcp", "127.0.0.1:37001")
		if err != nil {
			fmt.Println("Connection failed:", err)
			return
		}
		defer conn.Close()
	


		_, myErr := conn.Write(result)
		if myErr != nil {
			fmt.Println("Error writing:", myErr)
			return
		}

		fmt.Println("[WitcherConnect] Geralt Naked command was sent succesfully!")
}


// func blockAction(w http.ResponseWriter, r *http.Request) {
// 	command := "blockCrossbow()"
// 	fmt.Println("Request came through");
// 	time.Sleep(5 * time.Second) 

// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Unable to read request body", http.StatusBadRequest)
// 		return
// 	}

// 	var blockRequestBody BlockRequestBody
// 	err = json.Unmarshal(body, &blockRequestBody)
// 		if err != nil {
// 		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
// 		return
// 	}

// 	if (blockRequestBody.BlockedAction == "Signs") {
// 		command := "blockSigns()"

// 		result := Execute(command)
// 		// Attempt to connect to the specified IP and port
// 		conn, err := net.Dial("tcp", "127.0.0.1:37001")
// 		if err != nil {
// 			fmt.Println("Connection failed:", err)
// 			return
// 		}
// 		defer conn.Close()
	
// 		fmt.Println("Connection successful!")

		
// 	}

// 	// preCommand := "SpawnSpecialCompanion(\"Triss\", \"test\")"
// 	// command := strings.Replace(preCommand, "test", companionReqBody.Username, 1)
// 	fmt.Println(command)
// 	result := Execute(command)
// 		// Attempt to connect to the specified IP and port
// 		conn, err := net.Dial("tcp", "127.0.0.1:37001")
// 		if err != nil {
// 			fmt.Println("Connection failed:", err)
// 			return
// 		}
// 		defer conn.Close()
	
// 		fmt.Println("Connection successful!")

		
		
// 		_, myErr := conn.Write(result)
// 		if myErr != nil {
// 			fmt.Println("Error writing:", myErr)
// 			return
// 		}


// }


func getValidCompanion(companionName string) string {
	validCompanions := []string{"Triss", "Yennefer", "Ciri", "Vesemir", "Labert"}

	for _, companion := range validCompanions {
		if (strings.EqualFold(companionName, companion)) {
			return companion;
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return validCompanions[r.Intn(len(validCompanions))]
}


func spawnCompanion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[WitcherConnect] A request to spawn a companion came through!")
	time.Sleep(3 * time.Second) 

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var companionReqBody CompanionRequestBody
	err = json.Unmarshal(body, &companionReqBody)
		if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	companionName := getValidCompanion(companionReqBody.CompanionName)

	displayName, err := getDisplayName(companionReqBody.Username);
	
	if err != nil {
		fmt.Println("Error fetching display name:", err)
	} else {
		fmt.Println("Display name:", displayName)
	}

	command := fmt.Sprintf("SpawnSpecialCompanion(\"%s\", \"%s\")", companionName, displayName)
	fmt.Println(command)
	result := Execute(command)
		conn, err := net.Dial("tcp", "127.0.0.1:37001")
		if err != nil {
			fmt.Println("Connection failed:", err)
			return
		}
		defer conn.Close()
	
		fmt.Println("Connection successful!")

		
		
		_, myErr := conn.Write(result)
		if myErr != nil {
			fmt.Println("Error writing:", myErr)
			return
		}
	
		fmt.Println("[WitcherConnect] We placed a companion in the world!")
	
		time.AfterFunc(10*time.Minute, func() {
			DespawnCompanion(companionName)
		})
}


func DespawnCompanion(companionName string) {
	fmt.Println("[WitcherConnect] A request to despawn a companion came through!")
	time.Sleep(3 * time.Second) 

	command := "RemoveSpecialCompanions()"
	fmt.Println(command)
	result := Execute(command)
		conn, err := net.Dial("tcp", "127.0.0.1:37001")
		if err != nil {
			fmt.Println("Connection failed:", err)
			return
		}
		defer conn.Close()
	
		fmt.Println("Connection successful!")

		_, myErr := conn.Write(result)
		if myErr != nil {
			fmt.Println("Error writing:", myErr)
			return
		}
	
		fmt.Println("[WitcherConnect] We removed a companion from the world!")


}

func main() {
	fmt.Println("Welcome to WitcherHub! None of this probably works so enjoy!");

	
	http.HandleFunc("/companion", spawnCompanion)
	http.HandleFunc("/map", placeMapPin)
	// http.HandleFunc("/blocksigns", blockAllSigns)
	http.HandleFunc("/monsters", spawnMonsters)
	http.HandleFunc("/nakey", nakedGeralt)
	// http.HandleFunc("/test", testPause)

	// makeGeraltNaked()
	err := http.ListenAndServe(":8081", nil) 
	if err != nil {
		fmt.Println("Error goes here")
	}
}