package main

import (
	"sort"

	"github.com/GeoNet/delta/meta"
)

func (n *Network) EnviroSensor(set *meta.Set, enviro string, label string) error {

	net, ok := set.Network(enviro)
	if !ok {
		return nil
	}

	valid := make(map[string]interface{})
	for _, f := range set.Features() {
		valid[f.Station] = true
	}

	for _, stn := range set.Stations() {
		if _, ok := valid[stn.Code]; !ok {
			continue
		}

		if stn.Network != enviro {
			continue
		}

		var sites []Site
		for _, site := range set.Sites() {
			if site.Station != stn.Code {
				continue
			}

			var sensors []Sensor
			for _, v := range set.InstalledSensors() {
				if v.Station != site.Station {
					continue
				}
				if v.Location != site.Location {
					continue
				}

				sensors = append(sensors, Sensor{
					Type:  label,
					Make:  v.Make,
					Model: v.Model,

					Azimuth:  v.Azimuth,
					Method:   v.Method,
					Dip:      v.Dip,
					Vertical: v.Vertical,
					North:    v.North,
					East:     v.East,

					StartDate: v.Start,
					EndDate:   v.End,
				})
			}

			sort.Slice(sensors, func(i, j int) bool {
				return sensors[i].Less(sensors[j])
			})

			sites = append(sites, Site{
				Code: site.Location,

				Latitude:  site.Latitude,
				Longitude: site.Longitude,
				Elevation: site.Elevation,
				Datum:     site.Datum,
				Survey:    site.Survey,

				StartDate: site.Start,
				EndDate:   site.End,

				Sensors: sensors,
			})
		}

		sort.Slice(sites, func(i, j int) bool {
			return sites[i].Less(sites[j])
		})

		n.Stations = append(n.Stations, Station{
			Code:        stn.Code,
			Network:     net.External,
			Name:        stn.Name,
			Description: net.Description,

			Latitude:  stn.Latitude,
			Longitude: stn.Longitude,
			Elevation: stn.Elevation,
			Depth:     stn.Depth,
			Datum:     stn.Datum,

			StartDate: stn.Start,
			EndDate:   stn.End,

			Sites: sites,
		})
	}

	return nil
}
