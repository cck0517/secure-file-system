package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	"bytes"
	_ "encoding/hex"
	"errors"
	_ "strconv"
	"strings"
	"testing"

	// "sync"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"
	// "github.com/google/uuid"
	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const wrongPassword = "wrongPassword"
const emptyString = ""
const contentOne = "Content ONE "
const contentTwo = "Content TWO "
const contentThree = "Content THREE"
const contentFour = "Content FOUR"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	var aliceCopy *client.User
	var doris *client.User
	var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User
	var bobPhone *client.User
	var bobLaptop *client.User
	var bobDesktop *client.User
	// var datastore map[uuid.UUID][]byte
	var datastoreCopy map[uuid.UUID][]byte
	var diff map[uuid.UUID][]byte
	var data []byte
	var aliceData []byte
	var bobData []byte

	var err error
	var invite uuid.UUID

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	aliceFile2 := "aliceFile2.txt"
	aliceFile3 := "aliceFile3.txt"
	aliceFile4 := "aliceFile4.txt"
	bobFile := "bobFile.txt"
	bobFile2 := "bobFile2.txt"
	bobFile3 := "bobFile3.txt"
	bobFile4 := "bobFile4.txt"
	charlesFile := "charlesFile.txt"
	charlesFile1Copy := "charlesFile1Copy.txt"
	charlesFile2 := "charlesFile2.txt"
	charlesFile3 := "charlesFile3.txt"
	charlesFile4 := "charlesFile4.txt"
	dorisFile := "dorisFile.txt"
	dorisFile2 := "dorisFile2.txt"
	dorisFile3 := "dorisFile3.txt"
	dorisFile4 := "dorisFile4.txt"
	eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and charles.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and charles accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	userlib.DebugMsg("\n Starting Custome Tests.")
	Describe("Custom Tests", func() {
		measureBandwidth := func(probe func()) (bandwidth int) {
			before := userlib.DatastoreGetBandwidth()
			probe()
			after := userlib.DatastoreGetBandwidth()
			return after - before
		}

		Specify("Custom Test: Testing InitUser/GetUser Edge Cases", func() {

			// datastore = userlib.DatastoreGetMap()

			// get user before initialize
			userlib.DebugMsg("1.1 Get user before initialize")
			alice, err = client.GetUser("alice", defaultPassword)
			Expect(err).NotTo(BeNil())

			// 1.1 Initializing user with empty username.
			userlib.DebugMsg("1.1 Initializing user with empty username or empty password.")
			alice, err = client.InitUser(emptyString, defaultPassword)
			Expect(err).To(MatchError(errors.New("Username cannot be empty!")))
			Expect(alice).To(BeNil())

			alice, err = client.InitUser("alice", emptyString)
			Expect(err).To(MatchError(errors.New("Password cannot be empty!")))
			Expect(alice).To(BeNil())

			// 1.2 Alice login without initialize
			userlib.DebugMsg("1.2 Alice login without initialize")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(MatchError(errors.New("Username does not exist or Userinfo has been tempered!")))
			Expect(aliceLaptop).To(BeNil())

			// 1.3	Initializing user with same username
			userlib.DebugMsg("1.3 Initializing user Alice with same username.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			// userlib.DebugMsg("Initializing user Alice Copy.")
			aliceCopy, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(MatchError(errors.New("Username already exists or data has been tempered!")))
			Expect(aliceCopy).To(BeNil())

			// 1.4 Get user with wrong username and wrong password
			userlib.DebugMsg("1.4 Alice login wrong username")
			aliceLaptop, err = client.GetUser("wrongUserName", defaultPassword)
			Expect(err).To(MatchError(errors.New("Username does not exist or Userinfo has been tempered!")))
			Expect(aliceLaptop).To(BeNil())

			aliceLaptop, err = client.GetUser("alice", wrongPassword)
			Expect(err).To(MatchError(errors.New("couldn't unmarshal user struct!")))
			Expect(aliceLaptop).To(BeNil())

			// 1.5 Alice login from two laptop and desktop
			userlib.DebugMsg("1.5 Alice login from two laptop and desktop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(aliceLaptop).ToNot(BeNil())
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(aliceDesktop).ToNot(BeNil())

			// 1.6 get user that is not initialzied
			userlib.DebugMsg("1.6 get user that is not initialzied")
			bobLaptop, err = client.GetUser("bob", defaultPassword)
			Expect(err).To(MatchError(errors.New("Username does not exist or Userinfo has been tempered!")))

			// initizalie Bob, so that there are two elements in the map
			bob, err = client.InitUser("bob", defaultPassword)
			bob.StoreFile(bobFile, []byte(contentOne))
			userlib.DebugMsg("1.7 Bob try to get alice's file")
			bobData, err = bob.LoadFile(aliceFile)
			Expect(err).NotTo(BeNil())
			Expect(bobData).To(BeNil())

			// 1.7 DataStore is tampered
			// userlib.DebugMsg("1.8 DataStore is tampered")
			// // change the value in the first postion in the map to random bytes
			// for key, _ := range userlib.DatastoreGetMap() {
			// 	userlib.DatastoreSet(key, []byte("randomBytes"))
			// 	break
			// }
			// aliceLaptopTampered, err := client.GetUser("alice", defaultPassword)
			// Expect(err).NotTo(BeNil())
			// Expect(aliceLaptopTampered).To(BeNil())
			// bobLaptopTampered, err := client.GetUser("bob", defaultPassword)
			// Expect(err).NotTo(BeNil())
			// Expect(bobLaptopTampered).To(BeNil())

			// // load file should be error
			// bobData, err = bob.LoadFile(bobFile)
			// Expect(err).NotTo(BeNil())
			// Expect(bobData).To(BeNil())
		})

		Specify("Custom Test: Tamper with InitUser/GetUser", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			userlib.DebugMsg("Copy the datastore before initializing bob")
			datastoreCopy = copyDataStore(userlib.DatastoreGetMap())
			bob, err = client.InitUser("bob", defaultPassword)
			bobPhone, err = client.GetUser("bob", defaultPassword)
			Expect(err).To(BeNil())
			Expect(bobPhone).NotTo(BeNil())
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())
			userlib.DebugMsg("Tamper the datastore before initializing bob")
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("bob cannot login after the datastore is tampered")
			bobPhone, err = client.GetUser("bob", defaultPassword)
			Expect(err).NotTo(BeNil())
			Expect(bobPhone).To(BeNil())
			userlib.DebugMsg("alice can still login")
			alice, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			Expect(alice).NotTo(BeNil())
		})

		Specify("Custom Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("2.1 store a file that does not exist.")
			err = aliceDesktop.StoreFile("notExistFile", []byte(contentOne))
			Expect(err).NotTo(BeNil())

			// 2.1 Load a file that does not exist.
			userlib.DebugMsg("2.1 Load a file that does not exist.")
			aliceLoadFile, err := aliceDesktop.LoadFile("notExistFile")
			Expect(err).To(MatchError(errors.New("file not found in user's namespace")))
			Expect(aliceLoadFile).To(BeNil())

			// 2.2 Store a file and load
			userlib.DebugMsg("2.2 Store a file and load")
			aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			dataDesktop, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(dataDesktop).To(Equal([]byte(contentOne)))

			// 2.3 store a file that exists alread
			userlib.DebugMsg("2.3 store a file that exists alread")
			aliceDesktop.StoreFile(aliceFile, []byte(contentTwo))
			dataDesktop, err = aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(dataDesktop).To(Equal([]byte(contentTwo)))

			// 2.4 load a file from different device
			userlib.DebugMsg("2.4 load a file from different device")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
			dataLaptop, err := aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(dataLaptop).To(Equal([]byte(contentTwo)))

			// 2.5 append to file that does not exist
			userlib.DebugMsg("2.5 append to file that does not exist")
			err = aliceDesktop.AppendToFile("nonexistFile", []byte(contentThree))
			Expect(err).To(MatchError(errors.New("couldn't getmeta")))

			// 2.6 append to file a empty string
			userlib.DebugMsg("2.6 append to file a empty string")
			err = aliceDesktop.AppendToFile(aliceFile, []byte(""))
			Expect(err).NotTo(BeNil())

			// 2.6 append to file from a different device and load
			userlib.DebugMsg("2.6 append to file from a different device and load")
			err = alicePhone.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())
			dataLaptop, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(dataLaptop).To(Equal([]byte(contentTwo + contentThree)))

			// // 2.7 tamper data in the datastore
			// userlib.DebugMsg("2.7 tamper data in the datastore")
			// for key, _ := range userlib.DatastoreGetMap() {
			// 	userlib.DatastoreSet(key, []byte("randomBytes"))
			// }
			// err = aliceDesktop.StoreFile(aliceFile2, []byte(contentOne))
			// Expect(err).To(MatchError(errors.New("Couldn't verify and decrypt!")))
			// dataLaptopTampered, err := aliceLaptop.LoadFile(aliceFile)
			// Expect(err).To(MatchError(errors.New("Couldn't verify and decrypt!")))
			// Expect(dataLaptopTampered).To(BeNil())
			// err = aliceDesktop.AppendToFile(aliceFile, []byte(contentTwo))
			// Expect(err).To(MatchError(errors.New("couldn't getmeta")))
		})
		userlib.DebugMsg("Custome Test: Testing baisc Create/Accept Invitation")
		Specify("Custome Test: Testing baisc Create/Accept Invitation", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			userlib.DebugMsg("3.1 createInvitaion that the filename doesn't exist")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(MatchError(errors.New("file not found in user's namespace")))
			Expect(invite).To(Equal(uuid.Nil))

			userlib.DebugMsg("3.2 recipient doesn't exist")
			invite, err = alice.CreateInvitation(aliceFile, "notExistRecipient")
			Expect(err).NotTo(BeNil())
			Expect(invite).To(Equal(uuid.Nil))

			alice.StoreFile(aliceFile, []byte(contentOne))
			userlib.DebugMsg("3.3 createInvitaion to recipient more than one times")
			invite, err = alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			alice.StoreFile(aliceFile2, []byte(contentOne))
			userlib.DebugMsg("Create invitaion many times to bob before bob accept")
			invite, err = alice.CreateInvitation(aliceFile2, "bob")
			Expect(err).To(BeNil())
			invite, err = alice.CreateInvitation(aliceFile2, "bob")
			Expect(err).To(BeNil())
			invite, err = alice.CreateInvitation(aliceFile2, "bob")
			Expect(err).To(BeNil())
			invite, err = alice.CreateInvitation(aliceFile2, "bob")
			Expect(err).To(BeNil())
			// bob accepts
			err = bob.AcceptInvitation("alice", invite, bobFile2)
			Expect(err).To(BeNil())

		})

		userlib.DebugMsg("Custom Test: User Store/Load/Append fuzz case")
		Specify("Custom Test: User Store/Load/Append fuzz case", func() {
			alice, err = client.InitUser(strings.Repeat("alice", 1000000), defaultPassword)
			Expect(err).To(BeNil())

			alice, err = client.InitUser("alice", strings.Repeat("alice", 10000000))
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob try to login Alice's account")
			bobPhone, err = client.GetUser("bob", defaultPassword)
			bobLaptop, err = client.GetUser("bob", defaultPassword)
			bobDesktop, err = client.GetUser("bob", defaultPassword)
			userlib.DebugMsg("Bob's phone store bobfile")
			err = bobPhone.StoreFile(bobFile, []byte(contentOne))
			Expect(err).To(BeNil())
			userlib.DebugMsg("Bob's desktop try to access the bobfile")
			bobData, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(bobData).To(Equal([]byte(contentOne)))
			userlib.DebugMsg("Bob's desktop append content two to bobfile")
			err = bobDesktop.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())
			userlib.DebugMsg("Bob's laptop append content three to bobfile")
			err = bobLaptop.AppendToFile(bobFile, []byte(contentThree))
			Expect(err).To(BeNil())
			userlib.DebugMsg("Bob's phone load bobfile")
			bobData, err = bobPhone.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(bobData).To(Equal([]byte(contentOne + contentTwo + contentThree)))
			userlib.DebugMsg("Bob's laptop load bobfile")
			bobData, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(bobData).To(Equal([]byte(contentOne + contentTwo + contentThree)))
			userlib.DebugMsg("Bob's desktop load bobfile")
			bobData, err = bobDesktop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(bobData).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		userlib.DebugMsg("Custom Test: Tamper with Store")
		Specify("Custom Test: Tamper with Store", func() {
			// get the copy
			datastoreCopy = userlib.DatastoreGetMap()
			// initialize user
			alice, err = client.InitUser("alice", defaultPassword)
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())
			userlib.DebugMsg("Tamper the datastore before store")
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("bob cannot StoreFile after the datastore is tampered")
			// store file
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).NotTo(BeNil())
			// loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).NotTo(BeNil())
			Expect(aliceData).To(BeNil())
		})
		userlib.DebugMsg("Custom Test: Tamper with Load")
		Specify("Custom Test: Tamper with Store", func() {
			// initialize user
			alice, err = client.InitUser("alice", defaultPassword)
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			// get the copy
			datastoreCopy = userlib.DatastoreGetMap()

			// store file
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).NotTo(BeNil())
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())
			userlib.DebugMsg("Tamper the datastore before store")
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("bob cannot LoadFile after the datastore is tampered")
			// loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).NotTo(BeNil())
			Expect(aliceData).To(BeNil())
		})

		userlib.DebugMsg("Custom Test: Tamper with AppendFile")
		Specify("Custom Test: Tamper with Store", func() {
			// initialize user
			alice, err = client.InitUser("alice", defaultPassword)
			aliceDesktop, err = client.GetUser("alice", defaultPassword)
			// get the copy
			datastoreCopy = userlib.DatastoreGetMap()

			// store file
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())

			// loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(aliceData).To(Equal([]byte(contentOne)))

			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("bob cannot appendFile after the datastore is tampered")
			// append file
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).NotTo(BeNil())
			// loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).NotTo(BeNil())
			Expect(aliceData).To(BeNil())
		})

		userlib.DebugMsg("Custom Test: Tamper after AppendFile")
		Specify("Custom Test: Tamper with Store", func() {
			// initialize user
			alice, err = client.InitUser("alice", defaultPassword)
			aliceDesktop, err = client.GetUser("alice", defaultPassword)

			// store file
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(aliceData).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("bob cannot LoadFile after the datastore is tampered")
			// append file
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).NotTo(BeNil())
			//loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(aliceData).To(Equal([]byte(contentOne + contentTwo)))
			// get the copy
			datastoreCopy = userlib.DatastoreGetMap()
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).NotTo(BeNil())
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("Cannot append to file when the datastore is tampered")
			// appendfile
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentFour))
			Expect(err).NotTo(BeNil())

			// loadfile
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).NotTo(BeNil())
			Expect(aliceData).To(BeNil())
		})

		Specify("Custome Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			charles, err = client.InitUser("charles", defaultPassword)
			eve, err = client.InitUser("eve", defaultPassword)
			doris, err = client.InitUser("doris", defaultPassword)
			// frank, err := client.InitUser("frank", defaultPassword)

			// Alice store file
			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			// Alice shares the file to Bob
			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			// Bob Accecpts
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
			// Bob shares the file to Charles
			userlib.DebugMsg("Bob creating invite for Charles for file %s, and charles accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("3.1 Charles load the file before accept the invitation")
			data, err := charles.LoadFile(charlesFile)
			Expect(err).To(MatchError(errors.New("file not found in user's namespace")))
			Expect(data).To(BeNil())

			// Charles Accepts
			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("3.2 Bob share the same file to Charles again")
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())
			err = charles.AcceptInvitation("bob", invite, charlesFile1Copy)
			Expect(err).To(BeNil())

			userlib.DebugMsg("3.3 Alice shares a nonexist file")
			invite, err = alice.CreateInvitation("nonexistfile", "charles")
			Expect(err).To(MatchError(errors.New("file not found in user's namespace")))
			Expect(invite).To(Equal(uuid.Nil))

			userlib.DebugMsg("3.3 Alice shares file to a nonexist user")
			invite, err = alice.CreateInvitation(aliceFile, "nonexistuser")
			Expect(err).NotTo(BeNil())
			Expect(invite).To(Equal(uuid.Nil))

			userlib.DebugMsg("Alice store file2 and share it to Charles")
			alice.StoreFile(aliceFile2, []byte(contentOne))
			invite, err = alice.CreateInvitation(aliceFile2, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Charles accepts the invitaion of File2 from Alice")
			err = charles.AcceptInvitation("alice", invite, charlesFile2)
			Expect(err).To(BeNil())

			// Charles share the file to Doris
			userlib.DebugMsg("Charles share file1 to Doris")
			invite, err = charles.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Doris accepts the invitaion of File1 from charles")
			err = doris.AcceptInvitation("charles", invite, dorisFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice share file to Eve")
			invite, err = alice.CreateInvitation(aliceFile, "eve")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Eve accepts the invitaion of File from Alice")
			err = eve.AcceptInvitation("alice", invite, eveFile)
			Expect(err).To(BeNil())

			//  		Alice(file1)
			// 			|   \
			// 		  Eve 	bob
			// 				  \
			// 				  charles
			// 				  |
			// 				doris

			userlib.DebugMsg("3.4 Eve try to share the same file to Alice again")
			invite, err = eve.CreateInvitation(eveFile, "alice")
			Expect(err).To(BeNil())
			err = alice.AcceptInvitation("eve", invite, aliceFile3)
			Expect(err).To(BeNil())
			userlib.DebugMsg("AliceFile3 should be the same as AliceFile")
			data, err = alice.LoadFile(aliceFile)
			data3, err := alice.LoadFile(aliceFile3)
			Expect(err).To(BeNil())
			Expect(data).To(Equal(data3))

			userlib.DebugMsg("3.4 Alice revoke access of a person that she didn't share with")
			err = alice.RevokeAccess(aliceFile, "frank")
			Expect(err).To(MatchError(errors.New("recipient not found in invitation list")))

			userlib.DebugMsg("3.4.2 alice try to revoke bob's file")
			err = alice.RevokeAccess(bobFile, "bob")
			Expect(err).NotTo(BeNil())

			userlib.DebugMsg("3.5 Bob try to revoke file that he is not the owner")
			err = bob.RevokeAccess(bobFile, "eve")
			Expect(err).To(MatchError(errors.New("only the owner can revoke access")))

			userlib.DebugMsg("3.6 Alice try to revoke a user that Alice didn't directly shared the file with")
			err = alice.RevokeAccess(aliceFile, "charles")
			Expect(err).To(MatchError(errors.New("recipient not found in invitation list")))

			userlib.DebugMsg("3.7 Alice revoke access of Bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Check if Eve can still access the file")
			data, err = eve.LoadFile(eveFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("3.8 Check if Charles, Doris, Bob cannot access the file")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())

			userlib.DebugMsg("3.9 Check that CharlesCopy cannot access the file")
			data, err = charles.LoadFile(charlesFile1Copy)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())

			//  		Alice(file1)
			// 			|
			// 		   Eve

		})

		Specify("Custome Test: Nested file share", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			charles, err = client.InitUser("charles", defaultPassword)
			doris, err = client.InitUser("doris", defaultPassword)
			// frank, err := client.InitUser("frank", defaultPassword)

			// Alice store file
			userlib.DebugMsg("Alice storing four files")
			alice.StoreFile(aliceFile, []byte(contentOne))
			alice.StoreFile(aliceFile2, []byte(contentTwo))
			alice.StoreFile(aliceFile3, []byte(contentThree))
			alice.StoreFile(aliceFile4, []byte(contentFour))
			// Alice shares the file to Bob
			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			invite2, err := alice.CreateInvitation(aliceFile2, "bob")
			invite3, err := alice.CreateInvitation(aliceFile3, "bob")
			invite4, err := alice.CreateInvitation(aliceFile4, "bob")
			Expect(err).To(BeNil())
			// Bob Accecpts
			err = bob.AcceptInvitation("alice", invite, bobFile)
			err = bob.AcceptInvitation("alice", invite2, bobFile2)
			err = bob.AcceptInvitation("alice", invite3, bobFile3)
			err = bob.AcceptInvitation("alice", invite4, bobFile4)
			Expect(err).To(BeNil())
			// Bob shares the file to Charles
			userlib.DebugMsg("Bob creating invite for Charles for file %s, and charles accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "Charles")
			invite2, err = bob.CreateInvitation(bobFile2, "Charles")
			invite3, err = bob.CreateInvitation(bobFile3, "Charles")
			invite4, err = bob.CreateInvitation(bobFile4, "Charles")
			Expect(err).To(BeNil())
			// Charles Accepts
			err = charles.AcceptInvitation("bob", invite, charlesFile)
			err = charles.AcceptInvitation("bob", invite2, charlesFile2)
			err = charles.AcceptInvitation("bob", invite3, charlesFile3)
			err = charles.AcceptInvitation("bob", invite4, charlesFile4)
			Expect(err).To(BeNil())

			// Charles share the file to Doris
			userlib.DebugMsg("Charles share file1 to Doris")
			invite, err = charles.CreateInvitation(charlesFile, "doris")
			invite2, err = charles.CreateInvitation(charlesFile2, "doris")
			invite3, err = charles.CreateInvitation(charlesFile3, "doris")
			invite4, err = charles.CreateInvitation(charlesFile4, "doris")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Doris accepts the invitaion of File1 from Charles")
			err = doris.AcceptInvitation("charles", invite, dorisFile)
			err = doris.AcceptInvitation("charles", invite2, dorisFile2)
			err = doris.AcceptInvitation("charles", invite3, dorisFile3)
			err = doris.AcceptInvitation("charles", invite4, dorisFile4)
			Expect(err).To(BeNil())
			//  		Alice
			// 			   \
			// 		  		bob
			// 				  \
			// 				  charles
			// 				  |
			// 				doris

			userlib.DebugMsg("Alice revoek the access of file1 for Bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Bob, Charles, and Doris cannot access fiel1")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())
			data, err = doris.LoadFile(dorisFile)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())

			userlib.DebugMsg("Bob revoek the access of file2 for Charles")
			err = bob.RevokeAccess(bobFile2, "charles")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Alice can access file2")
			data, err = alice.LoadFile(aliceFile2)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)))

			userlib.DebugMsg("Charles, Doris cannot access file2")
			data, err = charles.LoadFile(charlesFile2)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())
			data, err = doris.LoadFile(dorisFile2)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))
			Expect(data).To(BeNil())

			userlib.DebugMsg("Charles revoek the access of file3 for Doris")
			err = charles.RevokeAccess(charlesFile3, "doris")
			Expect(err).To(BeNil())
			userlib.DebugMsg("Alice and bob can access file2")
			data, err = alice.LoadFile(aliceFile2)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)))
			data, err = bob.LoadFile(bobFile2)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentTwo)))

			userlib.DebugMsg("Doris cannot access file3")
			data, err = doris.LoadFile(dorisFile3)
			Expect(err).To(MatchError(errors.New("Couldn't retrieve invitation data!")))

		})

		Specify("Custome Test: Nested file share", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			charles, err = client.InitUser("charles", defaultPassword)
			eve, err = client.InitUser("eve", defaultPassword)
			doris, err = client.InitUser("doris", defaultPassword)
			// frank, err := client.InitUser("frank", defaultPassword)

			// Alice store file
			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			// Alice shares the file to Bob
			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			// Bob Accecpts
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
			// Bob shares the file to Charles
			userlib.DebugMsg("Bob creating invite for Charles for file %s, and charles accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())
			// Charles Accepts
			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			// Charles share the file to Doris
			userlib.DebugMsg("Charles share file1 to Doris")
			invite, err = charles.CreateInvitation(charlesFile, "doris")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Doris accepts the invitaion of File1 from charles")
			err = doris.AcceptInvitation("charles", invite, dorisFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice share file to Eve")
			invite, err = alice.CreateInvitation(aliceFile, "eve")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Eve accepts the invitaion of File from Alice")
			err = eve.AcceptInvitation("alice", invite, eveFile)
			Expect(err).To(BeNil())

			//  		Alice(file1)
			// 			|   \
			// 		  Eve 	bob
			// 				  \
			// 				  charles
			// 				  |
			// 				doris

			userlib.DebugMsg("3.11 revoke acess form the middle of the tree")
			err = bob.RevokeAccess(bobFile, "charles")
			Expect(err).To(BeNil())

			userlib.DebugMsg("check if charles can access the file")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).NotTo(BeNil())
			Expect(data).To(BeNil())

			userlib.DebugMsg("check if doris can access the file")
			data, err = doris.LoadFile(dorisFile)
			Expect(err).NotTo(BeNil())
			Expect(data).To(BeNil())

			userlib.DebugMsg("check if eve can access the file")
			data, err = eve.LoadFile(eveFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("check if alice can access the file")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("check if bob can access the file")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("check if bob laptop can access the file")
			data, err = bobLaptop.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

		})

		Specify("Custome Test: 4. More edge cases.", func() {

			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			eve, err = client.InitUser("eve", defaultPassword)
			// Alice store file
			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			// Alice shares the file to Bob
			userlib.DebugMsg("Alice creating invite for Bob")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			// Bob Accecpts
			userlib.DebugMsg("Bob accept the invitation")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Alice revokes the invitation from Bob")
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("4.3 check the validality of invitation after it is revoked")
			// bob accept the invitation again
			userlib.DebugMsg("Bob tries to accept the invitation again")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).NotTo(BeNil())

			userlib.DebugMsg("4.5 Alice try to revoke again")
			err = alice.RevokeAccess(aliceFile, "eve")
			Expect(err).To(MatchError(errors.New("recipient not found in invitation list")))

		})

		Specify("Custome Test: 4. Bob accept invitaion to a existing file.", func() {
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)

			// Alice store file
			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))
			bob.StoreFile(bobFile, []byte(contentOne))
			// Alice shares the file to Bob
			userlib.DebugMsg("Alice creating invite for Bob")
			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			// Bob Accecpts
			userlib.DebugMsg("Bob accept the invitation")
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).NotTo(BeNil())
		})

		Specify("5 Custome Test: Tamper the datastore before creating the invitation", func() {
			// Tamper Part
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			// copy datastore
			datastoreCopy = copyDataStore(userlib.DatastoreGetMap())
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// get the difference between storing the file
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())
			userlib.DebugMsg("5.1 Tamper the datastore before creating the invitation")
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("Alice cannot access the file")
			aliceData, err = alice.LoadFile(aliceFile)
			Expect(err).NotTo(BeNil())
			Expect(aliceData).To(BeNil())

			userlib.DebugMsg("CreateInvitation, AcceptInvitation, and RevokeAccess should all not work")
			invite, err = alice.CreateInvitation(aliceFile, "bob")
			Expect(err).NotTo(BeNil())
			Expect(invite).To(Equal(uuid.Nil))
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).NotTo(BeNil())
			err = alice.RevokeAccess(aliceFile, "eve")
			Expect(err).NotTo(BeNil())

		})

		Specify("Custome Test: more edge cases in invitation", func() {
			// initialize alice, bob, eve
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			eve, err = client.InitUser("eve", defaultPassword)
			// alice stores file
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			// alice creates the invitation to bob
			invite, err = alice.CreateInvitation(aliceFile, "bob")
			// bob accept file to a nonexist file
			userlib.DebugMsg("Bob accept file and store it into a nonexistFile")
			err = bob.AcceptInvitation("alice", invite, "nonexistFile")
			Expect(err).NotTo(BeNil())

			// alice shares another file to bob, but eve tries to accept the invitation
			userlib.DebugMsg("Alice shares another file to bob, but eve tries to accept the invitation")
			err = alice.StoreFile(aliceFile2, []byte(contentOne))
			invite, err = alice.CreateInvitation(aliceFile2, "bob")
			Expect(err).To(BeNil())
			// eve accept the invitation
			userlib.DebugMsg("Eve accept the invitation")
			err = eve.AcceptInvitation("alice", invite, eveFile)
			Expect(err).NotTo(BeNil())

		})

		Specify("5.2 Custome Test: Tamper the datastore after creating the invitation", func() {
			// Tamper Part
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			// get the copy of datastore before creating the invitation
			datastoreCopy = copyDataStore(userlib.DatastoreGetMap())
			invite, err = alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			// get the difference between creating the invitation
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())
			userlib.DebugMsg("5.2 Tamper the datastore before accepting the invitation")
			// Tamper the invitation in the datastore
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("AcceptInvitation and RevokeAccess should all not work")

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
			err = alice.RevokeAccess(aliceFile, "eve")
			Expect(err).NotTo(BeNil())

		})
		Specify("5.3 Custome Test: Tamper the datastore after accepting the invitation", func() {
			// Tamper Part
			alice, err = client.InitUser("alice", defaultPassword)
			bob, err = client.InitUser("bob", defaultPassword)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			invite, err = alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())
			// get the copy of datastore before accepting the invitation
			datastoreCopy = copyDataStore(userlib.DatastoreGetMap())
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())
			userlib.DebugMsg("Bob can access the file")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))
			// get the difference between accepting the invitation
			diff = diffMaps(datastoreCopy, userlib.DatastoreGetMap())

			userlib.DebugMsg("5.3 Tamper the datastore after accepting the invitation")
			// Tamper the invitation in the datastore
			for key, _ := range diff {
				userlib.DatastoreSet(key, []byte("randomBytes"))
			}
			userlib.DebugMsg("RevokeAccess should all not work")
			userlib.DebugMsg("Bob cannot access the file")
			data, err = bob.LoadFile(bobFile)
			Expect(err).NotTo(BeNil())
			Expect(data).To(BeNil())

			err = alice.RevokeAccess(aliceFile, "eve")
			Expect(err).NotTo(BeNil())

		})

		Specify("Custome Test: Test Bandwidth", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err := client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			// Preparing contents
			smallContent := []byte(contentOne)
			largeContent := userlib.RandomBytes(1024 * 1024)
			appendContent := []byte(contentTwo)

			// Store small content file
			userlib.DebugMsg("Alice stores small content file")
			err = alice.StoreFile("smallFile.txt", smallContent)
			Expect(err).To(BeNil())

			// Store large content file
			userlib.DebugMsg("Alice stores large content file")
			err = alice.StoreFile("largeFile.txt", largeContent)
			Expect(err).To(BeNil())

			// Measure bandwidth for appending to the small content file
			bandwidthSmallFile := measureBandwidth(func() {
				userlib.DebugMsg("Alice appends to small content file")
				err = alice.AppendToFile("smallFile.txt", appendContent)
				Expect(err).To(BeNil())
			})

			// Measure bandwidth for appending to the large content file
			bandwidthLargeFile := measureBandwidth(func() {
				userlib.DebugMsg("Alice appends to large content file")
				err = alice.AppendToFile("largeFile.txt", appendContent)
				Expect(err).To(BeNil())
			})

			// Comparing bandwidth usage to check efficiency
			userlib.DebugMsg("Comparing bandwidth for appending operations")
			Expect(bandwidthSmallFile).To(BeNumerically("<=", bandwidthLargeFile))

		})

	})

})

func copyDataStore(datastore map[uuid.UUID][]byte) map[uuid.UUID][]byte {
	datastoreCopy := make(map[uuid.UUID][]byte)
	for key, value := range datastore {
		datastoreCopy[key] = value
	}
	return datastoreCopy
}

// compare the difference between two map, and store the key and values into diff
func diffMaps(datastore map[uuid.UUID][]byte, datastoreCopy map[uuid.UUID][]byte) map[uuid.UUID][]byte {
	diff := make(map[uuid.UUID][]byte)
	for key, value := range datastoreCopy {
		if !bytes.Equal(value, datastore[key]) {
			diff[key] = value
		}
	}
	return diff
}
