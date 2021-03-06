package doorman

import (
	"context"

	authv1beta1 "github.com/OchiengEd/doorman-operator/pkg/apis/auth/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDoorman) reconcileDoormanService(cr *authv1beta1.Doorman) error {
	svc := createDoormanService(cr)
	found := r.isObjectFound(types.NamespacedName{Name: svc.Name, Namespace: cr.Namespace}, svc)
	if found {
		return nil
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.scheme); err != nil {
		return err
	}

	return r.client.Create(context.TODO(), svc)
}

func createDoormanService(cr *authv1beta1.Doorman) *corev1.Service {
	labels := map[string]string{
		"app.kubernetes.io/component":  "authentication",
		"app.kubernetes.io/managed-by": "doorman-operator",
		"app.kubernetes.io/name":       cr.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "api",
					Protocol: corev1.ProtocolTCP,
					Port:     5000,
					TargetPort: intstr.IntOrString{
						IntVal: 5000,
					},
				},
			},
			Selector: map[string]string{
				"app": cr.Name,
			},
		},
	}
}

func (r *ReconcileDoorman) reconcileDatabaseService(cr *authv1beta1.Doorman) error {
	svc := createDatabaseService(cr)
	found := r.isObjectFound(types.NamespacedName{Name: svc.Name, Namespace: cr.Namespace}, svc)
	if found {
		return nil
	}

	if err := controllerutil.SetControllerReference(cr, svc, r.scheme); err != nil {
		return err
	}

	return r.client.Create(context.TODO(), svc)
}

// Create headless service for the database
func createDatabaseService(cr *authv1beta1.Doorman) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name + "-database",
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-database",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: corev1.ClusterIPNone,
			Selector: map[string]string{
				"app": cr.Name + "-database",
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}
