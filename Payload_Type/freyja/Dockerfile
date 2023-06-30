FROM itsafeaturemythic/mythic_go_base_image:latest

WORKDIR /Mythic/

COPY [".", "."]

RUN make build

CMD make run
