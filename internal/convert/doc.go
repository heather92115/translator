// Package convert provides a collection of utility functions for converting between internal models and GraphQL types within
// a Go application. These functions facilitate the mapping of Go structs to GraphQL objects and vice versa, which is crucial
// for integrating Go backends with GraphQL APIs. The package covers conversions for various types,
// including but not limited to:
//
// DateTime Conversion: Functions to convert between Go's time.Time type and GraphQL's DateTime string format,
// enabling accurate representation of time data in GraphQL requests and responses.
//
// Audit Records: Mapping functions for converting audit logs, encapsulating data such as changes and states before
// and after modifications, from internal representations to GraphQL models, and vice versa. This includes handling of
// complex structures like diffs and timestamps.
//
// Vocabulary Data: Conversions between internal vocabulary structures and their GraphQL counterparts, ensuring
// that language learning or processing apps can effectively communicate vocab-related data through GraphQL.
//
// This package serves as a bridge for data transformation, ensuring that the Go backend can seamlessly operate with
// GraphQL frontends or clients, thus supporting a wide range of applications from web services to data processing
// tools in domains requiring detailed audit logs or language processing capabilities.
//
// Key Features
// Ease of Integration: Simplifies the integration of Go backends with GraphQL by providing ready-to-use conversion logic.
// Time Handling: Includes utilities for handling time data, crucial for applications that need to work with dates and times.
// Audit and Vocabulary Support: Offers specialized conversions for audit logs and vocabulary data, catering to
// applications with needs for tracking changes and language data management.
//
// This package is designed to be used in conjunction with a GraphQL server implementation in Go, enhancing data handling
// and API communication by bridging the gap between Go's type system and GraphQL's type definitions.
package convert
