package job

type Repository interface {
	Create(job Job) error
	Get(id string) (Job, error)
	List() ([]Job, error)
}
