

run:
	cp components.go otelcol-dev/
	go run ./otelcol-dev --config tailtracer/config.yaml

build:
	builder --config=otelcol-builder.yaml

rpm:
	rpmbuild -bb rpm.spec --define "_sourcedir ${PWD}"

