package flows

// func (f FlowRun) Handler(event interface{}) error {
// 	filterMatch := false
// 	for

// 	if err := mapstructure.Decode(event, &data); err != nil {
// 		return err
// 	}
// 	// Business Logic
// 	var rhelVersion, baseosURL, appstreamURL, imageID string
// 	var codereadyContainersMessage bool = false
// 	for _, product := range data.Artifact.Products {
// 		if product.Name == "rhel" {
// 			rhelVersion = product.Id
// 			baseosURL, appstreamURL = getRepositoriesURLs(product.Repos)
// 			logging.Debugf("Got repos baseos: %s, appstream %s", baseosURL, appstreamURL)
// 			imageID = product.Image
// 		}
// 		if product.Name == "codeready_containers" {
// 			codereadyContainersMessage = true
// 		}
// 	}
// 	// Filtering this will be improved in future versions
// 	if len(rhelVersion) > 0 && codereadyContainersMessage {
// 		name, xunitURL, duration, resultStatus, err :=
// 			interopPipelineRHEL.Run(rhelVersion, baseosURL, appstreamURL, imageID)
// 		if err != nil {
// 			logging.Error(err)
// 		}
// 		// We will take info from status to send back the results
// 		response := buildResponse(name, xunitURL, duration, resultStatus, &data)
// 		return umb.Send(topicTestComplete, response)
// 	}
// 	return nil
// }
