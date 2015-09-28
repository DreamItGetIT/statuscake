// Package statuscacke implements a client for statuscake.com API.
//
//  // list all `Tests`
//  c := statuscake.New(username, apikey)
//
//  tests, err := c.Tests().All()
//  if err != nil {
//    log.Fatal(err)
//  }
//
//  // delete a `Test`
//  err = c.Tests().Delete(TestID)
//
//  // create a test
//  t := &statuscake.Test{
//    WebsiteName: "Foo",
//    WebsiteURL:  "htto://example.com",
//    ... other required args...
//  }
//
//  if err = t.Validate(); err != nil {
//    log.Fatal(err)
//  }
//
//  t2 := c.Tests().Put(t)
//  fmt.Printf("New Test created with id: %d\n", t2.TestID)
package statuscake
