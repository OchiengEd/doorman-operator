package doorman

import (
	"context"

	authv1beta1 "github.com/OchiengEd/doorman-operator/pkg/apis/auth/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDoorman) reconcileStatefulSet(cr *authv1beta1.Doorman) error {
	statefulset := createDoormanDB(cr)
	found := r.isObjectFound(types.NamespacedName{Name: statefulset.Name, Namespace: cr.Namespace}, statefulset)
	if found {
		return nil
	}

	if err := controllerutil.SetControllerReference(cr, statefulset, r.scheme); err != nil {
		return err
	}

	return r.client.Create(context.TODO(), statefulset)
}

func createDoormanDB(cr *authv1beta1.Doorman) *appsv1.StatefulSet {
	labels := map[string]string{
		"app": cr.Name + "-database",
	}

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-database",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       cr.Name,
				"app.kubernetes.io/managed-by": "doorman-operator",
				"app.kubernetes.io/component":  "database",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			// Replicas:    &cr.Spec.Replicas,
			Replicas:    &cr.Spec.Storage.Size,
			ServiceName: cr.Name + "-database",
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						// MySQL container
						{
							Name:  "mariadb",
							Image: "mariadb",
							Env: []corev1.EnvVar{
								{
									Name: "MYSQL_ROOT_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: cr.Name + "-db",
											},
											Key: "root",
										},
									},
								},
								{
									Name: "MYSQL_USER",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: cr.Name + "-db",
											},
											Key: "username",
										},
									},
								},
								{
									Name: "MYSQL_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: cr.Name + "-db",
											},
											Key: "password",
										},
									},
								},
								{
									Name: "MYSQL_DATABASE",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: cr.Name + "-db",
											},
											Key: "database",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "db",
									ContainerPort: 3306,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
}
