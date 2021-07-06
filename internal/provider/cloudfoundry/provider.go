package cloudfoundry

import "github.com/joscha-alisch/dyve/pkg/provider/sdk"

func NewAppProvider(db Database) sdk.AppProvider {
	return &provider{
		db: db,
	}
}

type provider struct {
	db Database
}

func (p *provider) ListApps() ([]sdk.App, error) {
	cfApps, err := p.db.ListApps()
	if err != nil {
		return nil, err
	}

	var res []sdk.App
	for _, app := range cfApps {
		res = append(res, sdk.App{
			Id:   app.Guid,
			Name: app.Name,
		})
	}
	return res, nil
}

func (p *provider) GetApp(id string) (sdk.App, error) {
	app, err := p.db.GetApp(id)
	if err != nil {
		return sdk.App{}, err
	}

	return sdk.App{
		Id:   app.Guid,
		Name: app.Name,
	}, nil
}

func (p *provider) Search(term string, limit int) ([]sdk.AppSearchResult, error) {
	panic("implement me")
}