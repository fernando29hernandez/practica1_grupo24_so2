#include <linux/fs.h>
#include <linux/hugetlb.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/mm.h>
#include <linux/mman.h>
#include <linux/mmzone.h>
#include <linux/proc_fs.h>
#include <linux/quicklist.h>
#include <linux/seq_file.h>
#include <linux/swap.h>
#include <linux/vmstat.h>
#include <linux/atomic.h>
#include <asm/page.h>
#include <asm/pgtable.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/uaccess.h>
#include <linux/slab.h>
#include <linux/sysinfo.h>
#include <linux/slab.h>
void __attribute__((weak)) arch_report_meminfo(struct seq_file *m)
{
}
static int meminfo_proc_show(struct seq_file *m, void *v)
{
        struct sysinfo i;
/*
 * display in kilobytes.
 */
#define K(x) ((x) << (PAGE_SHIFT - 10))
si_meminfo(&i);
        seq_printf(m,
                "MemTotal,%8lu\n"
                "MemFree,%8lu\n"
                ,
                K(i.totalram),
                K(i.freeram)
                );

        return 0;
#undef K
}

static int meminfo_proc_open(struct inode *inode, struct file *file)
{
        return single_open(file, meminfo_proc_show, NULL);
}

static const struct file_operations meminfo_proc_fops = {
        .open           = meminfo_proc_open,
        .read           = seq_read,
        .llseek         = seq_lseek,
        .release        = single_release,
};

static int __init proc_meminfo_init(void)
{
    printk(KERN_INFO "“Hola mundo, somos el grupo 24 y este es el monitor de memoria\n");
        proc_create("mem_grupo24", 0, NULL, &meminfo_proc_fops);
        return 0;
}
static void __exit proc_meminfo_exit(void) {
 printk(KERN_INFO "Sayonara mundo, somos el grupo 24 y este fue el monitor de memoria\n");
 remove_proc_entry("mem_grupo24",NULL);
}
module_init(proc_meminfo_init);
module_exit(proc_meminfo_exit);

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Fernando Hernandez");
MODULE_DESCRIPTION("CPU");
MODULE_VERSION("0.01");
