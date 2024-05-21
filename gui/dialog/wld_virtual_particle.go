package dialog

import (
	"github.com/xackery/quail-gui/slog"
	"github.com/xackery/quail/vwld"
	"github.com/xackery/wlk/cpl"
	"github.com/xackery/wlk/walk"
)

func virtualParticlePage(data *vwld.VWld, page *cpl.TabPage) error {

	particles := []string{}
	for _, particle := range data.Particles {
		particles = append(particles, particle.Tag)
	}
	onParticleNew := func() {
		slog.Println("new particle")
	}
	onParticleEdit := func() {
		slog.Println("edit particle")
	}
	onParticleDelete := func() {
		slog.Println("delete particle")
	}

	var cmbParticle *walk.ComboBox
	defaultParticle := ""
	if len(particles) > 0 {
		defaultParticle = particles[0]
	}

	particleGroup := cpl.Composite{
		Layout: cpl.HBox{},
	}

	particleGroup.Children = append(particleGroup.Children, cpl.GroupBox{
		Title:  "Particles (BlitSpriteDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbParticle,
				Editable: false,
				Model:    particles,
				Value:    defaultParticle,
			},
			cpl.PushButton{Text: "Add", OnClicked: onParticleNew},
			cpl.PushButton{Text: "Edit", OnClicked: onParticleEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onParticleDelete},
		},
	})

	particleInstances := []string{}
	for _, particleInstance := range data.ParticleInstances {
		particleInstances = append(particleInstances, particleInstance.Tag)
	}
	onParticleInstanceNew := func() {
		slog.Println("new particleInstance")
	}
	onParticleInstanceEdit := func() {
		slog.Println("edit particleInstance")
	}
	onParticleInstanceDelete := func() {
		slog.Println("delete particleInstance")
	}

	var cmbParticleInstance *walk.ComboBox
	defaultParticleInstance := ""
	if len(particleInstances) > 0 {
		defaultParticleInstance = particleInstances[0]
	}

	particleGroup.Children = append(particleGroup.Children, cpl.GroupBox{
		Title:  "ParticleInstances (ParticleCloudDef)",
		Layout: cpl.HBox{},
		Children: []cpl.Widget{
			cpl.ComboBox{
				AssignTo: &cmbParticleInstance,
				Editable: false,
				Model:    particleInstances,
				Value:    defaultParticleInstance,
			},
			cpl.PushButton{Text: "Add", OnClicked: onParticleInstanceNew},
			cpl.PushButton{Text: "Edit", OnClicked: onParticleInstanceEdit},
			cpl.PushButton{Text: "Delete", OnClicked: onParticleInstanceDelete},
		},
	})

	page.Title = "Particle"
	page.Layout = cpl.VBox{}
	page.Children = []cpl.Widget{particleGroup}
	return nil
}
