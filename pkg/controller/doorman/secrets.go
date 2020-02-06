package doorman

import (
	"context"
	"math/rand"
	"time"

	authv1beta1 "github.com/OchiengEd/doorman-operator/pkg/apis/auth/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDoorman) reconcileDoormanSecrets(cr *authv1beta1.Doorman) error {
	databaseSecret := createDatabaseSecret(cr)
	foundDB := r.isObjectFound(types.NamespacedName{Name: databaseSecret.Name, Namespace: cr.Namespace}, databaseSecret)
	if foundDB {
		return nil
	}

	certSecret := createCertificateSecret(cr)
	foundCert := r.isObjectFound(types.NamespacedName{Name: certSecret.Name, Namespace: cr.Namespace}, certSecret)
	if foundCert {
		return nil
	}

	if err := controllerutil.SetControllerReference(cr, databaseSecret, r.scheme); err != nil {
		return err
	}

	if err := controllerutil.SetControllerReference(cr, certSecret, r.scheme); err != nil {
		return err
	}

	if err := r.client.Create(context.TODO(), certSecret); err != nil {
		return err
	}

	return r.client.Create(context.TODO(), databaseSecret)
}

func createDatabaseSecret(cr *authv1beta1.Doorman) *corev1.Secret {
	dbSecret := genericOpaqueSecret()
	dbSecret.ObjectMeta.Name = cr.Name + "-db"
	dbSecret.ObjectMeta.Namespace = cr.Namespace
	dbSecret.StringData = map[string]string{
		"database": cr.Spec.Storage.DatabaseName,
		"username": cr.Spec.Storage.Username,
		"root":     generatePassword(cr.Spec.Storage.PasswordLength),
		"password": generatePassword(cr.Spec.Storage.PasswordLength),
	}
	return dbSecret
}

func genericOpaqueSecret() *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "somesecret",
			Namespace: "default",
		},
		StringData: map[string]string{},
		Type:       corev1.SecretTypeOpaque,
	}
}

func createCertificateSecret(cr *authv1beta1.Doorman) *corev1.Secret {
	certAndKey, err := generateCertificateForKey()
	if err != nil {
		log.Error(err, "Error acquiring cert and key")
	}

	keyPairSecret := genericOpaqueSecret()
	keyPairSecret.Name = cr.Name + "-cert"
	keyPairSecret.Namespace = cr.Namespace
	keyPairSecret.StringData = map[string]string{
		corev1.TLSCertKey:       string(certAndKey.Certificate),
		corev1.TLSPrivateKeyKey: string(certAndKey.RSAPrivateKey),
	}
	keyPairSecret.Type = corev1.SecretTypeTLS

	return keyPairSecret
}

func generatePassword(size int) string {
	charset := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567890$&%*()-_=?/#@!+<>,")
	password := make([]byte, size)
	seed := rand.New(rand.NewSource(time.Now().Unix()))

	for i := range password {
		password[i] = charset[seed.Intn(len(charset))]
	}

	return string(password)
}
