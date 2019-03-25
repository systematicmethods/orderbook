A small package to 

*  Generate Time based Version 1 UUIDs from a time.Time.
   *  Needed for back testing as "github.com/google/uuid" does not provide
this feature.
*  UUIDComparator convenience func to compare order of Time based UUID
   * works with generated and time.Now based UUIDs
