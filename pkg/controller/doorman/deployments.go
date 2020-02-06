package doorman

import (
	"context"

	authv1beta1 "github.com/OchiengEd/doorman-operator/pkg/apis/auth/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileDoorman) reconcileDeployments(cr *authv1beta1.Doorman) error {
	deploy := newDoormanDeployment(cr)
	found := r.isObjectFound(types.NamespacedName{Name: deploy.Name, Namespace: cr.Namespace}, deploy)
	if found {
		return nil
	}

	if err := controllerutil.SetControllerReference(cr, deploy, r.scheme); err != nil {
		return err
	}

	return r.client.Create(context.TODO(), deploy)
}

func (r *ReconcileDoorman) reconcileServices(cr *authv1beta1.Doorman) error {
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
		"app": cr.Name,
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
			Selector: map[string]string{},
		},
	}
}

// newDoormanDeployment returns a Doorman resource with the same name/namespace as the cr
func newDoormanDeployment(cr *authv1beta1.Doorman) *appsv1.Deployment {
	certSecretName := cr.Name + "-cert"
	databaseSecretName := cr.Name + "-db"
	labels := map[string]string{
		"app": cr.Name,
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       cr.Name,
				"app.kubernetes.io/managed-by": "doorman-operator",
				"app.kubernetes.io/component":  "authentication",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{ //List of containers in pod
							Image: "quay.io/eochieng/doorman:v0.0.3",
							Name:  cr.Name,
							Env: []corev1.EnvVar{
								{
									Name: "DATABASE_USERNAME",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: databaseSecretName,
											},
											Key: "username",
										},
									},
								},
								{
									Name: "DATABASE_PASSWORD",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: databaseSecretName,
											},
											Key: "password",
										},
									},
								},
								{
									Name:  "DATABASE_HOST",
									Value: cr.Name + "-database",
								},
								{
									Name: "DATABASE",
									ValueFrom: &corev1.EnvVarSource{
										SecretKeyRef: &corev1.SecretKeySelector{
											LocalObjectReference: corev1.LocalObjectReference{
												Name: databaseSecretName,
											},
											Key: "database",
										},
									},
								},
							},
							Ports: []corev1.ContainerPort{
								{ //List of cont
									Name:          "auth",
									ContainerPort: 5000,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "certs",
									MountPath: "/etc/doorman/certs",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "certs",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: certSecretName,
								},
							},
						},
					},
				},
			},
		},
	}
}
