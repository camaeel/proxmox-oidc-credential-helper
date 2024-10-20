package proxmox

type Ticket struct {
	Data TicketData `json:"data"`
}

type TicketData struct {
	Ticket              string `json:"ticket"`
	CSRFPreventionToken string `json:"CSRFPreventionToken"`
	Username            string `json:"username"`
	Clustername         string `json:"clustername,omitempty"`
}
