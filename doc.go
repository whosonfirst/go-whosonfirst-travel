// Package travel provides methods for iterating through one or more Who's On First records and "traveling" through the other records that it has a relationship with. These include parent records, records that supersede or are superseded by a WOF record or all the pointers in a records hierarchy.
//
// Each step (or record) in a travel function is processed by a user-defined `TravelFunc`. This package includes a default callback function that simply prints some basic information about each record. More complex applications are outside the scope of this package.
package travel
