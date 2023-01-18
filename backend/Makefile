.PHONY: oas-code-gen-echo
oas-code-gen-echo:
	oapi-codegen --old-config-style -generate types -package echo presentation/oas.yml > presentation/echo/types.gen.go
	oapi-codegen --old-config-style -generate server -package echo presentation/oas.yml > presentation/echo/server.gen.go
