package main

import (
	"fmt"
	"log"

	ldap "github.com/go-ldap/ldap/v3"
)

func main() {
	// Connect to LDAP
	conn, err := ldap.DialURL("ldap://localhost:30000")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Bind with admin credentials
	err = conn.Bind("cn=admin,dc=devplatform,dc=local", "admin123")
	if err != nil {
		log.Fatalf("Failed to bind: %v", err)
	}

	fmt.Println("✓ Connected to LDAP as admin")

	// Create base OUs
	ous := []struct {
		dn   string
		ou   string
		desc string
	}{
		{"ou=users,dc=devplatform,dc=local", "users", "Users"},
		{"ou=groups,dc=devplatform,dc=local", "groups", "Groups"},
		{"ou=departments,dc=devplatform,dc=local", "departments", "Departments"},
	}

	for _, ou := range ous {
		addReq := ldap.NewAddRequest(ou.dn, nil)
		addReq.Attribute("objectClass", []string{"organizationalUnit"})
		addReq.Attribute("ou", []string{ou.ou})
		addReq.Attribute("description", []string{ou.desc})

		err = conn.Add(addReq)
		if err != nil {
			if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
				fmt.Printf("⚠ OU %s already exists\n", ou.ou)
			} else {
				log.Printf("Failed to create OU %s: %v", ou.ou, err)
			}
		} else {
			fmt.Printf("✓ Created OU: %s\n", ou.ou)
		}
	}

	// Create a test user
	userDN := "uid=john.doe,ou=users,dc=devplatform,dc=local"
	addReq := ldap.NewAddRequest(userDN, nil)
	addReq.Attribute("objectClass", []string{"inetOrgPerson", "posixAccount", "shadowAccount"})
	addReq.Attribute("uid", []string{"john.doe"})
	addReq.Attribute("cn", []string{"John Doe"})
	addReq.Attribute("sn", []string{"Doe"})
	addReq.Attribute("givenName", []string{"John"})
	addReq.Attribute("mail", []string{"john.doe@devplatform.local"})
	addReq.Attribute("uidNumber", []string{"10001"})
	addReq.Attribute("gidNumber", []string{"10001"})
	addReq.Attribute("homeDirectory", []string{"/home/john.doe"})
	addReq.Attribute("userPassword", []string{"password123"})

	err = conn.Add(addReq)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			fmt.Println("⚠ User john.doe already exists")
		} else {
			log.Fatalf("Failed to create user: %v", err)
		}
	} else {
		fmt.Println("✓ Created user: john.doe")
	}

	fmt.Println("\n✅ LDAP Initialization Complete!")
	fmt.Println("\nLogin Credentials:")
	fmt.Println("  Username: john.doe")
	fmt.Println("  Password: password123")
}
