package template

var templateStr = `
version: "3.5"

services:
    {{range .Nodes}}
    {{ .Name }}:
        image: {{ .Image }}
        container_name: {{ .ContainerName }}
        command: {{ .Command }}
        ports:
           - "{{ .PortP2P }}"
           - "{{ .PortHttp }}"
           - "{{ .PortWs }}"
        networks:
             qlcchain:
                ipv4_address: {{ .Ipv4Address }}
        volumes:
           - {{ .Volumes }}
        restart: unless-stopped
    {{end}}

    {{range .Ptms}}
    {{ .Name }}:
        image: {{ .Image }}
        container_name: {{ .ContainerName }}
        command: {{ .Command }}
        ports:
           - "{{ .Port1 }}"
           - "{{ .Port2 }}"
        networks:
             qlcchain:
                ipv4_address: {{ .Ipv4Address }}
        volumes:
           - {{ .Volumes }}
        restart: unless-stopped
    {{end}}

networks:
   qlcchain:
      ipam:
         config:
         - subnet: 10.0.0.0/16
`
