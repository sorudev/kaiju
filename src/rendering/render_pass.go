package rendering

import (
	"kaiju/klib"
	vk "kaiju/rendering/vulkan"
	"log/slog"
)

type RenderPassAttachmentImage struct {
	MipLevels      uint32
	LayerCount     uint32
	Tiling         string   // vk.ImageTiling
	Filter         string   // vk.Filter
	Usage          []string // vk.ImageUsageFlagBits
	MemoryProperty []string // vk.MemoryPropertyFlagBits
	Aspect         []string // vk.ImageAspectFlagBits
	Access         []string // vk.AccessFlagBits
}

type RenderPassAttachmentDescription struct {
	Format         string
	Samples        string
	LoadOp         string
	StoreOp        string
	StencilLoadOp  string
	StencilStoreOp string
	InitialLayout  string
	FinalLayout    string
	Image          RenderPassAttachmentImage
}

type RenderPassSubpassDescription struct {
	PipelineBindPoint         string
	ColorAttachmentReferences []RenderPassAttachmentReference
	InputAttachmentReferences []RenderPassAttachmentReference
	ResolveAttachments        []RenderPassAttachmentReference
	DepthStencilAttachment    []RenderPassAttachmentReference // 1 max
	PreserveAttachments       []uint32                        // TODO
}

type RenderPassSubpassDependency struct {
	SrcSubpass      uint32
	DstSubpass      uint32
	SrcStageMask    []string
	DstStageMask    []string
	SrcAccessMask   []string
	DstAccessMask   []string
	DependencyFlags []string
}

type RenderPassAttachmentReference struct {
	Attachment uint32
	Layout     string
}

func (ai *RenderPassAttachmentImage) ListTiling() []string {
	return klib.MapKeysSorted(StringVkImageTiling)
}

func (ai *RenderPassAttachmentImage) ListFilter() []string {
	return klib.MapKeysSorted(StringVkFilter)
}

func (ai *RenderPassAttachmentImage) ListUsage() []string {
	return klib.MapKeysSorted(StringVkImageUsageFlagBits)
}

func (ai *RenderPassAttachmentImage) ListMemoryProperty() []string {
	return klib.MapKeysSorted(StringVkMemoryPropertyFlagBits)
}

func (ai *RenderPassAttachmentImage) ListAspect() []string {
	return klib.MapKeysSorted(StringVkImageAspectFlagBits)
}

func (ai *RenderPassAttachmentImage) ListAccess() []string {
	return klib.MapKeysSorted(StringVkAccessFlagBits)
}

func (ai *RenderPassAttachmentImage) TilingToVK() vk.ImageTiling {
	return imageTilingToVK(ai.Tiling)
}

func (ai *RenderPassAttachmentImage) FilterToVK() vk.Filter {
	return filterToVK(ai.Filter)
}

func (ai *RenderPassAttachmentImage) UsageToVK() vk.ImageUsageFlags {
	return imageUsageFlagsToVK(ai.Usage)
}

func (ai *RenderPassAttachmentImage) MemoryPropertyToVK() vk.MemoryPropertyFlags {
	return memoryPropertyFlagsToVK(ai.MemoryProperty)
}

func (ai *RenderPassAttachmentImage) AspectToVK() vk.ImageAspectFlags {
	return imageAspectFlagsToVK(ai.Aspect)
}

func (ai *RenderPassAttachmentImage) AccessToVK() vk.AccessFlags {
	return accessFlagsToVK(ai.Access)
}

func (ad *RenderPassAttachmentDescription) ListFormat() []string {
	return klib.MapKeysSorted(StringVkFormat)
}

func (ad *RenderPassAttachmentDescription) ListSamples() []string {
	return klib.MapKeysSorted(StringVkSampleCountFlagBits)
}

func (ad *RenderPassAttachmentDescription) ListLoadOp() []string {
	return klib.MapKeysSorted(StringVkAttachmentLoadOp)
}

func (ad *RenderPassAttachmentDescription) ListStoreOp() []string {
	return klib.MapKeysSorted(StringVkAttachmentStoreOp)
}

func (ad *RenderPassAttachmentDescription) ListStencilLoadOp() []string {
	return klib.MapKeysSorted(StringVkAttachmentLoadOp)
}

func (ad *RenderPassAttachmentDescription) ListStencilStoreOp() []string {
	return klib.MapKeysSorted(StringVkAttachmentStoreOp)
}

func (ad *RenderPassAttachmentDescription) ListInitialLayout() []string {
	return klib.MapKeysSorted(StringVkImageLayout)
}

func (ad *RenderPassAttachmentDescription) ListFinalLayout() []string {
	return klib.MapKeysSorted(StringVkImageLayout)
}

func (ad *RenderPassAttachmentDescription) FormatToVK(vr *Vulkan) vk.Format {
	return formatToVK(ad.Format, vr)
}

func (ad *RenderPassAttachmentDescription) SamplesToVK() vk.SampleCountFlagBits {
	return sampleCountToVK(ad.Samples)
}

func (ad *RenderPassAttachmentDescription) LoadOpToVK() vk.AttachmentLoadOp {
	return attachmentLoadOpToVK(ad.LoadOp)
}

func (ad *RenderPassAttachmentDescription) StoreOpToVK() vk.AttachmentStoreOp {
	return attachmentStoreOpToVK(ad.StoreOp)
}

func (ad *RenderPassAttachmentDescription) StencilLoadOpToVK() vk.AttachmentLoadOp {
	return attachmentLoadOpToVK(ad.StencilLoadOp)
}

func (ad *RenderPassAttachmentDescription) StencilStoreOpToVK() vk.AttachmentStoreOp {
	return attachmentStoreOpToVK(ad.StencilStoreOp)
}

func (ad *RenderPassAttachmentDescription) InitialLayoutToVK() vk.ImageLayout {
	return imageLayoutToVK(ad.InitialLayout)
}

func (ad *RenderPassAttachmentDescription) FinalLayoutToVK() vk.ImageLayout {
	return imageLayoutToVK(ad.FinalLayout)
}

func (ad *RenderPassAttachmentReference) ListLayout() []string {
	return klib.MapKeysSorted(StringVkImageLayout)
}

func (ad *RenderPassAttachmentReference) LayoutToVK() vk.ImageLayout {
	return imageLayoutToVK(ad.Layout)
}

func (ad *RenderPassSubpassDescription) ListPipelineBindPoint() []string {
	return klib.MapKeysSorted(StringVkPipelineBindPoint)
}

func (ad *RenderPassSubpassDescription) PipelineBindPointToVK() vk.PipelineBindPoint {
	return pipelineBindPointToVK(ad.PipelineBindPoint)
}

func (sd *RenderPassSubpassDependency) ListSrcStageMask() []string {
	return klib.MapKeysSorted(StringVkPipelineStageFlagBits)
}

func (sd *RenderPassSubpassDependency) ListDstStageMask() []string {
	return klib.MapKeysSorted(StringVkPipelineStageFlagBits)
}

func (sd *RenderPassSubpassDependency) ListSrcAccessMask() []string {
	return klib.MapKeysSorted(StringVkAccessFlagBits)
}

func (sd *RenderPassSubpassDependency) ListDstAccessMask() []string {
	return klib.MapKeysSorted(StringVkAccessFlagBits)
}

func (sd *RenderPassSubpassDependency) ListDependencyFlags() []string {
	return klib.MapKeysSorted(StringVkDependencyFlagBits)
}

func (sd *RenderPassSubpassDependency) SrcStageMaskToVK() vk.PipelineStageFlags {
	return pipelineStageFlagsToVK(sd.SrcStageMask)
}

func (sd *RenderPassSubpassDependency) DstStageMaskToVK() vk.PipelineStageFlags {
	return pipelineStageFlagsToVK(sd.DstStageMask)
}

func (sd *RenderPassSubpassDependency) SrcAccessMaskToVK() vk.AccessFlags {
	return accessFlagsToVK(sd.SrcAccessMask)
}

func (sd *RenderPassSubpassDependency) DstAccessMaskToVK() vk.AccessFlags {
	return accessFlagsToVK(sd.DstAccessMask)
}

func (sd *RenderPassSubpassDependency) DependencyFlagsToVK() vk.DependencyFlags {
	return dependencyFlagsToVK(sd.DependencyFlags)
}

type RenderPassData struct {
	Name                   string
	AttachmentDescriptions []RenderPassAttachmentDescription
	SubpassDescriptions    []RenderPassSubpassDescription
	SubpassDependencies    []RenderPassSubpassDependency
}

func (r *RenderPassData) ConstructRenderPass(renderer Renderer) (RenderPass, bool) {
	vr := renderer.(*Vulkan)
	textures := make([]TextureId, len(r.AttachmentDescriptions))
	{
		w := uint32(vr.swapChainExtent.Width)
		h := uint32(vr.swapChainExtent.Height)
		for i := range len(r.AttachmentDescriptions) {
			a := &r.AttachmentDescriptions[i]
			img := &a.Image
			success := vr.CreateImage(w, h, img.MipLevels, a.SamplesToVK(),
				a.FormatToVK(vr), img.TilingToVK(), img.UsageToVK(),
				img.MemoryPropertyToVK(), &textures[i], int(img.LayerCount))
			if !success {
				slog.Error("failed to create image for render pass attachment", "attachmentIndex", i)
				return RenderPass{}, false
			}
			success = vr.createImageView(&textures[i], img.AspectToVK())
			if !success {
				for j := range i + 1 {
					vr.textureIdFree(&textures[j])
				}
				slog.Error("failed to create image view for render pass attachment", "attachmentIndex", i)
				return RenderPass{}, false
			}
			success = vr.createTextureSampler(&textures[i].Sampler,
				img.MipLevels, img.FilterToVK())
			if !success {
				for j := range i + 1 {
					vr.textureIdFree(&textures[j])
				}
				slog.Error("failed to create image sampler for render pass attachment", "attachmentIndex", i)
				return RenderPass{}, false
			}
			success = vr.transitionImageLayout(&textures[i], a.InitialLayoutToVK(),
				img.AspectToVK(), img.AccessToVK(), vk.NullCommandBuffer)
			if !success {
				for j := range i + 1 {
					vr.textureIdFree(&textures[j])
				}
				slog.Error("failed to transition image layout for render pass attachment", "attachmentIndex", i)
				return RenderPass{}, false
			}
		}
	}
	attachments := make([]vk.AttachmentDescription, len(r.AttachmentDescriptions))
	for i := range r.AttachmentDescriptions {
		// TODO:  Flags
		attachments[i].Flags = 0
		attachments[i].Format = r.AttachmentDescriptions[i].FormatToVK(vr)
		attachments[i].Samples = r.AttachmentDescriptions[i].SamplesToVK()
		attachments[i].LoadOp = r.AttachmentDescriptions[i].LoadOpToVK()
		attachments[i].StoreOp = r.AttachmentDescriptions[i].StoreOpToVK()
		attachments[i].StencilLoadOp = r.AttachmentDescriptions[i].StencilLoadOpToVK()
		attachments[i].StencilStoreOp = r.AttachmentDescriptions[i].StencilStoreOpToVK()
		attachments[i].InitialLayout = r.AttachmentDescriptions[i].InitialLayoutToVK()
		attachments[i].FinalLayout = r.AttachmentDescriptions[i].FinalLayoutToVK()
	}
	color := make([][]vk.AttachmentReference, len(r.SubpassDescriptions))
	input := make([][]vk.AttachmentReference, len(r.SubpassDescriptions))
	preserve := make([][]uint32, len(r.SubpassDescriptions))
	depthStencil := make([][]vk.AttachmentReference, len(r.SubpassDescriptions))
	resolve := make([][]vk.AttachmentReference, len(r.SubpassDescriptions))
	for i := range r.SubpassDescriptions {
		c := r.SubpassDescriptions[i].ColorAttachmentReferences
		n := r.SubpassDescriptions[i].InputAttachmentReferences
		p := r.SubpassDescriptions[i].PreserveAttachments
		d := r.SubpassDescriptions[i].DepthStencilAttachment
		r := r.SubpassDescriptions[i].ResolveAttachments
		color[i] = make([]vk.AttachmentReference, len(c))
		input[i] = make([]vk.AttachmentReference, len(n))
		preserve[i] = make([]uint32, len(p))
		depthStencil[i] = make([]vk.AttachmentReference, len(d))
		resolve[i] = make([]vk.AttachmentReference, len(r))
		for j := range c {
			color[i][j].Attachment = c[j].Attachment
			color[i][j].Layout = c[j].LayoutToVK()
		}
		for j := range n {
			input[i][j].Attachment = n[j].Attachment
			input[i][j].Layout = n[j].LayoutToVK()
		}
		copy(p, preserve[i])
		for j := range depthStencil {
			depthStencil[i][j].Attachment = d[j].Attachment
			depthStencil[i][j].Layout = d[j].LayoutToVK()
		}
		for j := range resolve {
			resolve[i][j].Attachment = r[j].Attachment
			resolve[i][j].Layout = r[j].LayoutToVK()
		}
	}
	subpasses := make([]vk.SubpassDescription, len(r.SubpassDescriptions))
	for i := range r.SubpassDescriptions {
		// TODO:  Fill in the flags
		subpasses[i].Flags = 0
		subpasses[i].PipelineBindPoint = r.SubpassDescriptions[i].PipelineBindPointToVK()
		subpasses[i].ColorAttachmentCount = uint32(len(color))
		subpasses[i].InputAttachmentCount = uint32(len(input))
		subpasses[i].PreserveAttachmentCount = uint32(len(preserve))
		if len(color) > 0 {
			subpasses[i].PColorAttachments = &color[i][0]
		}
		if len(input) > 0 {
			subpasses[i].PInputAttachments = &input[i][0]
		}
		if len(preserve) > 0 {
			subpasses[i].PPreserveAttachments = &preserve[i][0]
		}
		if len(depthStencil) > 0 {
			subpasses[i].PDepthStencilAttachment = &depthStencil[i][0]
		}
		if len(resolve) > 0 {
			subpasses[i].PResolveAttachments = &resolve[i][0]
		}
	}
	selfDependencies := make([]vk.SubpassDependency, len(r.SubpassDependencies))
	for i := range r.SubpassDependencies {
		selfDependencies[i].SrcSubpass = r.SubpassDependencies[i].SrcSubpass
		selfDependencies[i].DstSubpass = r.SubpassDependencies[i].DstSubpass
		selfDependencies[i].SrcStageMask = r.SubpassDependencies[i].SrcStageMaskToVK()
		selfDependencies[i].DstStageMask = r.SubpassDependencies[i].DstStageMaskToVK()
		selfDependencies[i].SrcAccessMask = r.SubpassDependencies[i].SrcAccessMaskToVK()
		selfDependencies[i].DstAccessMask = r.SubpassDependencies[i].DstAccessMaskToVK()
		selfDependencies[i].DependencyFlags = r.SubpassDependencies[i].DependencyFlagsToVK()
	}
	pass, err := NewRenderPass(vr.device, &vr.dbg,
		attachments, subpasses, selfDependencies)
	if err != nil {
		slog.Error("failed to create the render pass", "error", err)
		return RenderPass{}, false
	}
	imageViews := make([]vk.ImageView, len(textures))
	for i := range textures {
		imageViews[i] = textures[i].View
	}
	err = pass.CreateFrameBuffer(vr, imageViews,
		textures[0].Width, textures[0].Height)
	if err != nil {
		slog.Error("failed to create the frame buffer for the render pass", "error", err)
		return RenderPass{}, false
	}
	return pass, true
}
