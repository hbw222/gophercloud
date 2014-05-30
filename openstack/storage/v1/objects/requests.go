package objects

import (
	"fmt"
	"github.com/racker/perigee"
	storage "github.com/rackspace/gophercloud/openstack/storage/v1"
	"github.com/rackspace/gophercloud/openstack/utils"
	"net/http"
)

// ListResult is a *http.Response that is returned from a call to the List function.
type ListResult *http.Response

// DownloadResult is a *http.Response that is returned from a call to the Download function.
type DownloadResult *http.Response

// GetResult is a *http.Response that is returned from a call to the Get function.
type GetResult *http.Response

// List is a function that retrieves all objects in a container. It also returns the details
// for the container. To extract only the object information or names, pass the ListResult
// response to the ExtractInfo or ExtractNames function, respectively.
func List(c *storage.Client, opts ListOpts) (ListResult, error) {
	contentType := ""

	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	query := utils.BuildQuery(opts.Params)

	if !opts.Full {
		contentType = "text/plain"
	}

	url := c.GetContainerURL(opts.Container) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: h,
		Accept:      contentType,
	})
	return &resp.HttpResponse, err
}

// Download is a function that retrieves the content and metadata for an object.
// To extract just the content, pass the DownloadResult response to the ExtractContent
// function.
func Download(c *storage.Client, opts DownloadOpts) (DownloadResult, error) {
	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	query := utils.BuildQuery(opts.Params)

	url := c.GetObjectURL(opts.Container, opts.Name) + query
	resp, err := perigee.Request("GET", url, perigee.Options{
		MoreHeaders: h,
	})
	return &resp.HttpResponse, err
}

// Create is a function that creates a new object or replaces an existing object.
func Create(c *storage.Client, opts CreateOpts) error {
	var reqBody []byte

	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	query := utils.BuildQuery(opts.Params)

	content := opts.Content
	if content != nil {
		reqBody = make([]byte, content.Len())
		_, err = content.Read(reqBody)
		if err != nil {
			return err
		}
	}

	url := c.GetObjectURL(opts.Container, opts.Name) + query
	_, err = perigee.Request("PUT", url, perigee.Options{
		ReqBody:     reqBody,
		MoreHeaders: h,
	})
	return err
}

// Copy is a function that copies one object to another.
func Copy(c *storage.Client, opts CopyOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	h["Destination"] = fmt.Sprintf("/%s/%s", opts.NewContainer, opts.NewName)

	url := c.GetObjectURL(opts.Container, opts.Name)
	_, err = perigee.Request("COPY", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}

// Delete is a function that deletes an object.
func Delete(c *storage.Client, opts DeleteOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	query := utils.BuildQuery(opts.Params)

	url := c.GetObjectURL(opts.Container, opts.Name) + query
	_, err = perigee.Request("DELETE", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}

// Get is a function that retrieves the metadata of an object. To extract just the custom
// metadata, pass the GetResult response to the ExtractMetadata function.
func Get(c *storage.Client, opts GetOpts) (GetResult, error) {
	h, err := c.GetHeaders()
	if err != nil {
		return nil, err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	url := c.GetObjectURL(opts.Container, opts.Name)
	resp, err := perigee.Request("HEAD", url, perigee.Options{
		MoreHeaders: h,
	})
	return &resp.HttpResponse, err
}

// Update is a function that creates, updates, or deletes an object's metadata.
func Update(c *storage.Client, opts UpdateOpts) error {
	h, err := c.GetHeaders()
	if err != nil {
		return err
	}

	for k, v := range opts.Headers {
		h[k] = v
	}

	for k, v := range opts.Metadata {
		h["X-Object-Meta-"+k] = v
	}

	url := c.GetObjectURL(opts.Container, opts.Name)
	_, err = perigee.Request("POST", url, perigee.Options{
		MoreHeaders: h,
	})
	return err
}
