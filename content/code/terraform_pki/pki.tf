# define a few values for my use case, replace it with your naming
locals {
  cn       = "rulz.xyz"
  org      = "Rulz Corp."
  validity = 87600 # This is approximatively 10 years
}

# To build a certificate you actually need a private key, so we are creating it here
resource "tls_private_key" "rulz_ca" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

# Here is our root certificate, using our private key previously created
resource "tls_self_signed_cert" "rulz_ca" {
  key_algorithm     = "ECDSA"
  private_key_pem   = tls_private_key.rulz_ca.private_key_pem  # we use our private key
  is_ca_certificate = true                                     # this is a certificate authority

  subject {
    common_name  = local.cn
    organization = local.org
  }

  validity_period_hours = local.validity

  allowed_uses = [
    "cert_signing",  # this is what will make this certificate able to sign child certificates
  ]
}


# We need to create a key for the intermediate cert
resource "tls_private_key" "rulz_int" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

# Once we have the key we can create a request
resource "tls_cert_request" "rulz_int" {
  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.rulz_int.private_key_pem

  # the wildcard here does not have an impact but will help remind that this
  # is an intermediate signing all *.rulz.xyz certificates
  subject {
    common_name  = "*.${local.cn}"
  }
}

# And sign the request with our root ca
resource "tls_locally_signed_cert" "rulz_int" {
  cert_request_pem   = tls_cert_request.rulz_int.cert_request_pem
  ca_key_algorithm   = "ECDSA"
  ca_private_key_pem = tls_private_key.rulz_ca.private_key_pem
  ca_cert_pem        = tls_self_signed_cert.rulz_ca.cert_pem
  is_ca_certificate  = true

  validity_period_hours = local.validity / 10 # one year of validity

  allowed_uses = [
    "cert_signing",
  ]
}

resource "tls_private_key" "rulz_child" {
  algorithm   = "ECDSA"
  ecdsa_curve = "P384"
}

resource "tls_cert_request" "rulz_child" {

  key_algorithm   = "ECDSA"
  private_key_pem = tls_private_key.rulz_child.private_key_pem
  dns_names       = [ "child.${local.cn}" ]

  subject {
    common_name = "child.${local.cn}"
  }
}

resource "tls_locally_signed_cert" "rulz_child" {
  cert_request_pem   = tls_cert_request.rulz_child.cert_request_pem
  ca_key_algorithm   = "ECDSA"
  ca_private_key_pem = tls_private_key.rulz_int.private_key_pem
  ca_cert_pem        = tls_locally_signed_cert.rulz_int.cert_pem

  validity_period_hours = local.validity / 40 # valid for ~ 3 months

  # this will only allow to be used as a server
  allowed_uses = [
    "server_auth",
  ]
}

locals {
  path = "${path.module}/ca_certificate.pem"
}

resource "local_file" "rulz_ca" {
  filename        = local.path
  content         = tls_self_signed_cert.rulz_ca.cert_pem
  file_permission = "0644"

  provisioner "local-exec" {  # I want to remove the certificate when I destroy the infra
    when    = destroy
    command = "rm ${self.filename}"
  }
}

resource "local_file" "rulz_child_key" {
  filename        = "${path.module}/child.key.pem"
  content         = tls_private_key.rulz_child.private_key_pem
  file_permission = "0600" # we want to protect the private key from others reading

  provisioner "local-exec" {
    when    = destroy
    command = "rm ${self.filename}"
  }
}

resource "local_file" "rulz_child_cert" {
  filename        = "${path.module}/child.crt.pem"
  file_permission = "0644"

  content = join("", [
    tls_locally_signed_cert.rulz_child.cert_pem,
    tls_locally_signed_cert.rulz_int.cert_pem,
    tls_self_signed_cert.rulz_ca.cert_pem,
  ])

  provisioner "local-exec" {
    when    = destroy
    command = "rm ${self.filename}"
  }
}
