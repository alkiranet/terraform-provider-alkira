# Terraform Provider for Alkira

* Website: http://www.alkira.com

The Terraform provider for Alkira is a Terraform plugin to enable full
lifecycle management of Alkira resources. The provider is maintained
internally by Alkira engineering team.

Currently, the provider is still in early development phase and being
actively worked on.

## Getting Started

To start using the provider, you will need Terraform 0.12.x.

Now, you could hover to the `/example` directory to explore some
examples and start playing with the provider.


## Development Requirements

-	[Terraform 0.12.29](https://releases.hashicorp.com/terraform/0.12.29/) 0.12.x
-	[Go](https://golang.org/doc/install) 1.16.x (to build the provider
     plugin on all architectures, especially Apple M1)

**NOTES:** Please don't use the latest Terraform 0.13 yet. The latest
release has some backward compatiblity break changes. We are waiting
for it to be more stablized.



