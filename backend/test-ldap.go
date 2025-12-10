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

	fmt.Println("✓ Connected to LDAP server on localhost:30000")

	// Bind with admin credentials
	err = conn.Bind("cn=admin,dc=devplatform,dc=local", "admin123")
	if err != nil {
		log.Fatalf("Failed to bind: %v", err)
	}

	fmt.Println("✓ Successfully authenticated as admin")

	// Search for base DN
	searchRequest := ldap.NewSearchRequest(
		"dc=devplatform,dc=local",
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=*)",
		[]string{"dn", "dc", "o"},
		nil,
	)

	result, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	fmt.Printf("✓ Search successful! Found %d entries\n", len(result.Entries))

	for _, entry := range result.Entries {
		fmt.Printf("  DN: %s\n", entry.DN)
		for _, attr := range entry.Attributes {
			fmt.Printf("    %s: %v\n", attr.Name, attr.Values)
		}
	}

	// Search for OUs
	searchRequest = ldap.NewSearchRequest(
		"dc=devplatform,dc=local",
		ldap.ScopeSingleLevel,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		"(objectClass=organizationalUnit)",
		[]string{"dn", "ou"},
		nil,
	)

	result, err = conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("OU search failed: %v", err)
	}

	fmt.Printf("\n✓ Found %d organizational units:\n", len(result.Entries))
	for _, entry := range result.Entries {
		fmt.Printf("  - %s\n", entry.DN)
	}

	fmt.Println("\n✅ LDAP connection test successful!")
	fmt.Println("OpenLDAP is ready to use on localhost:30000")
}
