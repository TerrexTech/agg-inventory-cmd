package inventory

// InternalError represents an error when something goes wrong, and its our fault.
const InternalError = 2

// DatabaseError is when some operation related to Database, such as insert or find,
// goes wrong and the task cannot proceed.
const DatabaseError = 3
