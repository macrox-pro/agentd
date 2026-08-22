// Package transport provides cross-platform listeners and dialers for the daemon IPC socket.
//
// Owns: Unix socket and Windows named pipe path, listen, dial.
// Must not: gRPC handlers, business logic.
//
// Entry: Listen, Dial, DefaultEndpoint.
package transport
