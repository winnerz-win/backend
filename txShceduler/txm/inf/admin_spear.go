package inf

const (
	AdminSpear = "spear"
)

var (
	is_admin_spear = false
)

func SetAdminSpear() { is_admin_spear = true }

func IsAdminSpear() bool { return is_admin_spear }
