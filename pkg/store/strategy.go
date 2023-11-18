package store

import (
	"context"
	"fmt"
	"github.com/myoperator/k8saggregatorapiserver/pkg/apis/myingress/v1beta1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/klog/v2"
	"k8s.io/kube-openapi/pkg/validation/strfmt"
	"k8s.io/kube-openapi/pkg/validation/validate"
)

// REST implements a RESTStorage for API services against etcd
type REST struct {
	*genericregistry.Store
}

func (*REST) ShortNames() []string {
	return []string{"mi"}
}

func RESTInPeace(storage rest.StandardStorage, err error) rest.StandardStorage {
	if err != nil {
		err = fmt.Errorf("unable to create REST storage for a resource due to %v, will die", err)
		panic(err)
	}
	return storage
}

// NewStrategy 构建资源对象的增删改查策略
func NewStrategy(typer runtime.ObjectTyper) MyIngressStrategy {
	return MyIngressStrategy{typer, names.SimpleNameGenerator}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {

	apiserver, ok := obj.(*v1beta1.MyIngress)
	if !ok {
		return nil, nil, fmt.Errorf(" object is not a MyIngress")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), nil
}

// MatchMyIngress 标签和字段匹配器
func MatchMyIngress(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *v1beta1.MyIngress) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type MyIngressStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// WarningsOnUpdate 更新时发出的警告
func (s MyIngressStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	//TODO implement me
	return []string{}
}

// WarningsOnCreate 创建时发出警告
func (s MyIngressStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	//TODO implement me
	return []string{}
}

func (MyIngressStrategy) NamespaceScoped() bool {
	return true
}

// PrepareForCreate 创建前调用 hook
func (MyIngressStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	klog.Info("PrepareForCreate method...")
}

// PrepareForUpdate 更新前调用 hook
func (MyIngressStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	klog.Info("PrepareForUpdate method...")
}

// Validate 字段验证相关
func (MyIngressStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	errList := field.ErrorList{}
	// 验证 spec 字段内容
	validatePath := field.NewPath("spec")

	spec := obj.(*v1beta1.MyIngress).Spec
	schema := spec.OpenAPIDefinition().Schema

	err := validate.AgainstSchema(&schema, spec, strfmt.Default)

	if err != nil {
		errList = append(errList, field.Invalid(validatePath, spec, err.Error()))
	}
	return field.ErrorList{}
}

// ValidateUpdate 更新验证方法
func (MyIngressStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (MyIngressStrategy) AllowCreateOnUpdate() bool {
	return true
}

func (MyIngressStrategy) AllowUnconditionalUpdate() bool {
	return true
}

func (MyIngressStrategy) Canonicalize(obj runtime.Object) {
}
