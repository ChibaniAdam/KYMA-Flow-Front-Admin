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
	addUserReq := ldap.NewAddRequest(userDN, nil)
	addUserReq.Attribute("objectClass", []string{"inetOrgPerson", "posixAccount", "shadowAccount"})
	addUserReq.Attribute("uid", []string{"john.doe"})
	addUserReq.Attribute("cn", []string{"John Doe"})
	addUserReq.Attribute("sn", []string{"Doe"})
	addUserReq.Attribute("givenName", []string{"John"})
	addUserReq.Attribute("mail", []string{"john.doe@devplatform.local"})
	addUserReq.Attribute("uidNumber", []string{"10001"})
	addUserReq.Attribute("gidNumber", []string{"10001"})
	addUserReq.Attribute("homeDirectory", []string{"/home/john.doe"})
	addUserReq.Attribute("userPassword", []string{"password123"})

	err = conn.Add(addUserReq)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			fmt.Println("⚠ User john.doe already exists")
		} else {
			log.Fatalf("Failed to create user: %v", err)
		}
	} else {
		fmt.Println("✓ Created user: john.doe")
	}

	// Create a test department
	deptDN := "ou=TestDept,ou=departments,dc=devplatform,dc=local"
	addDeptReq := ldap.NewAddRequest(deptDN, nil)
	addDeptReq.Attribute("objectClass", []string{"organizationalUnit"})
	addDeptReq.Attribute("ou", []string{"TestDept"})
	addDeptReq.Attribute("description", []string{"This is a test department"})
	// Optional fields from schema (manager, members, repositories) can be stored as extra attributes if LDAP schema allows
	// addDeptReq.Attribute("manager", []string{"john.doe"})
	// addDeptReq.Attribute("members", []string{"john.doe"})
	// addDeptReq.Attribute("repositories", []string{"repo1", "repo2"})

	err = conn.Add(addDeptReq)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			fmt.Println("⚠ Department TestDept already exists")
		} else {
			log.Fatalf("Failed to create department: %v", err)
		}
	} else {
		fmt.Println("✓ Created department: TestDept")
	}

	fmt.Println("\n✅ LDAP Initialization Complete!")
	fmt.Println("\nLogin Credentials:")
	fmt.Println("  Username: john.doe")
	fmt.Println("  Password: password123")
}
